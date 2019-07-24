package queue

import (
	"encoding/json"
	"github.com/innogames/slack-bot/bot/storage"
	"github.com/innogames/slack-bot/client"
	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
	"sync"
)

const (
	waitIcon   = "coffee"
	doneIcon   = "white_check_mark"
	storageKey = "fallback_queue"
)

var logger *logrus.Logger
var runningCommands = map[string]slack.MessageEvent{}
var mu sync.Mutex

// AddRunningCommand registers a long running command, e.g. a running Jenkins job or watching a pull request
// it's doing following magic:
// - when using "trigger job XXX" and "then reply done" it will execute the "reply done" when the running command was finished
// - when the server got restarted, the fallbackCommand gets executed
// todo add a kill chan to kill long running command via command
// todo improve locking
func AddRunningCommand(event slack.MessageEvent, fallbackCommand string) chan bool {
	var queueKey string

	if fallbackCommand != "" {
		event.Text = fallbackCommand

		queueKey = getKey(event)
		storage.Write(storageKey, queueKey, event)
	}

	if logger != nil {
		logger.Infof("add a blocking process: %s", event.Text)
	}
	mu.Lock()
	defer mu.Unlock()

	runningCommands[getKey(event)] = event

	finished := make(chan bool, 1)
	go func() {
		defer close(finished)

		// wait until blocking task is over
		<-finished

		mu.Lock()
		delete(runningCommands, getKey(event))
		mu.Unlock()

		if queueKey != "" {
			storage.Delete(storageKey, queueKey)
		}
	}()

	return finished
}

// IsBlock checks if there is a blocking command registered for this user/channel
func IsBlocked(event slack.MessageEvent) bool {
	_, ok := runningCommands[getKey(event)]

	return ok
}

func executeFallbackCommand(logger *logrus.Logger) {
	// todo check null
	res, _ := storage.ReadAll(storageKey)
	storage.Delete(storageKey, "")

	var event slack.MessageEvent
	for _, eventString := range res {
		if err := json.Unmarshal([]byte(eventString), &event); err != nil {
			logger.Errorf("[Queue] Not unmarshalable: %s", err)
			continue
		}

		logger.Infof("[Queue] Booted! I'll trigger this command now: `%s`", event.Text)
		client.InternalMessages <- event
	}
}

func getKey(event slack.MessageEvent) string {
	return event.User + event.Channel
}
