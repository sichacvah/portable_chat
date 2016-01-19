package model

import (
	"encoding/json"
	"io"

	"github.com/pborman/uuid"
)

const (
	CHANNEL_OPEN    = "O"
	CHANNEL_PRIVATE = "P"
	CHANNEL_DIRECT  = "D"
	DEFAULT_CHANNEL = "general"
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

func (c *Channel) IsValid() bool {
	return len(c.Name) > 0 && len(c.Type) > 0
}

func (c *Channel) PreSave() {
	c.Id = uuid.New()
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

func ChannelMapToJson(o map[string]*Channel) string {
	b, err := json.Marshal(o)
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}

func ChannelMapFromJson(data io.Reader) map[string]*Channel {
	decoder := json.NewDecoder(data)
	var channels map[string]*Channel
	err := decoder.Decode(&channels)
	if err == nil {
		return channels
	} else {
		return nil
	}
}
