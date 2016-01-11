package model

import (
	"encoding/json"
	"io"
)

const (
	ACTION_TYPING             = "typing"
	ACTION_POSTED             = "posted"
	ACTION_POST_EDITED        = "post_edited"
	ACTION_POST_DELETED       = "post_deleted"
	ACTION_CHANNEL_VIEWED     = "channel_viewed"
	ACTION_NEW_USER           = "new_user"
	ACTION_USER_ADDED         = "user_added"
	ACTION_USER_REMOVED       = "user_removed"
	ACTION_PREFERENCE_CHANGED = "preference_changed"
)

type Message struct {
	UserId    string            `json:"user_id"`
	ChannelId string            `json:"channel_id"`
	Action    string            `json:"action"`
	Props     map[string]string `json:"props"`
}

func (m *Message) Add(key string, value string) {
	m.Props[key] = value
}

func NewMessage(userId string, channelId string, action string) *Message {
	return &Message{UserId: userId, ChannelId: channelId, Action: action, Props: make(map[string]string)}
}

func (m *Message) ToJson() string {
	b, err := json.Marshal(m)
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}

func MessageFromJson(data io.Reader) *Message {
	decoder := json.NewDecoder(data)
	var m Message
	err := decoder.Decode(&m)
	if err == nil {
		return &m
	} else {
		return nil
	}
}
