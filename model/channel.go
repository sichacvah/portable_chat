package model

import (
	"encoding/json"
	"io"
)

type Channel struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// ToJson convert a User to a json string
func (c *Channel) ToJson() string {
	b, err := json.Marshal(c)
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}

// UserFromJson will decode the imput and return a user
func ChannelFromJson(data io.Reader) *Channel {
	decoder := json.NewDecoder(data)

	var channel Channel
	err := decoder.Decode(&channel)
	if err == nil {
		return &channel
	} else {
		return nil
	}
}

func ChannelMapToJson(u map[string]*Channel) string {
	b, err := json.Marshal(u)
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}

func ChannelMapFromJson(data io.Reader) map[string]*User {
	decoder := json.NewDecoder(data)
	var channels map[string]*Channel
	err := decoder.Decode(&channels)
	if err == nil {
		return channels
	} else {
		return nil
	}
}
