package bot

import (
	"github.com/innogames/slack-bot/bot/util"
	"github.com/innogames/slack-bot/client"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

// ListenForMessages is blocking method to handle new incoming events...from different sources
func (b *Bot) ListenForMessages(ctx *util.ServerContext) {
	ctx.RegisterChild()
	defer ctx.ChildDone()

	// listen for old/deprecated RTM connection
	// https://api.slack.com/rtm
	var rtmChan chan slack.RTMEvent
	if b.slackClient.RTM != nil {
		rtmChan = b.slackClient.RTM.IncomingEvents
	}

	// initialize Socket Mode:
	// https://api.slack.com/apis/connections/socket
	var socketChan chan socketmode.Event
	if b.slackClient.Socket != nil {
		go b.slackClient.Socket.Run()
		socketChan = b.slackClient.Socket.Events
	}

	for {
		select {
		case event := <-socketChan:
			// message from Socket Mode
			b.handleSocketModeEvent(event)
		case event := <-rtmChan:
			// message received from user via deprecated RTM API
			switch message := event.Data.(type) {
			case *slack.HelloEvent:
				log.Info("Hello, the RTM connection is ready!")
			case *slack.MessageEvent:
				b.HandleMessage(message)
			case *slack.RTMError, *slack.UnmarshallingErrorEvent, *slack.RateLimitEvent, *slack.ConnectionErrorEvent:
				log.Error(event)
			}
		case message := <-client.InternalMessages:
			// e.g. triggered by "delay" or "macro" command. They are still executed in original event context
			// -> will post in same channel as the user posted the original command
			message.InternalMessage = true
			go b.handleMessage(message, false)
		case <-ctx.Done():
			if err := b.DisconnectRTM(); err != nil {
				log.Error(err)
			}
			return
		}
	}
}

func (b *Bot) handleSocketModeEvent(event socketmode.Event) {
	if b.slackClient.Socket != nil {
		b.slackClient.Socket.Ack(*event.Request)
	}

	switch event.Type {
	case socketmode.EventTypeConnectionError, socketmode.EventTypeErrorBadMessage, socketmode.EventTypeErrorWriteFailed, socketmode.EventTypeIncomingError, socketmode.EventTypeInvalidAuth:
		log.Warnf("Socket Mode error: %s - %s", event.Type, event.Data)
	case socketmode.EventTypeEventsAPI:
		b.handleEvent(event.Data.(slackevents.EventsAPIEvent))
	case socketmode.EventTypeInteractive:
		b.handleInteraction(event.Data.(slack.InteractionCallback))
	case socketmode.EventTypeConnected, socketmode.EventTypeConnecting, socketmode.EventTypeHello, socketmode.EventTypeDisconnect, socketmode.EventTypeSlashCommand:
		// ignore
	default:
		log.Infof("Unexpected event type received: %s\n", event.Type)
	}
}
