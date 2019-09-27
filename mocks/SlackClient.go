// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import slack "github.com/nlopes/slack"

// SlackClient is an autogenerated mock type for the SlackClient type
type SlackClient struct {
	mock.Mock
}

// AddReaction provides a mock function with given fields: name, item
func (_m *SlackClient) AddReaction(name string, item slack.ItemRef) {
	_m.Called(name, item)
}

// RemoveReaction provides a mock function with given fields: name, item
func (_m *SlackClient) RemoveReaction(name string, item slack.ItemRef) {
	_m.Called(name, item)
}

// Reply provides a mock function with given fields: event, text
func (_m *SlackClient) Reply(event slack.MessageEvent, text string) {
	_m.Called(event, text)
}

// ReplyError provides a mock function with given fields: event, err
func (_m *SlackClient) ReplyError(event slack.MessageEvent, err error) {
	_m.Called(event, err)
}

// SendMessage provides a mock function with given fields: event, text, options
func (_m *SlackClient) SendMessage(event slack.MessageEvent, text string, options ...slack.MsgOption) string {
	_va := make([]interface{}, len(options))
	for _i := range options {
		_va[_i] = options[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, event, text)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 string
	if rf, ok := ret.Get(0).(func(slack.MessageEvent, string, ...slack.MsgOption) string); ok {
		r0 = rf(event, text, options...)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// SendToUser provides a mock function with given fields: user, text
func (_m *SlackClient) SendToUser(user string, text string) string {
	ret := _m.Called(user, text)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(user, text)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
