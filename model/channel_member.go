package model

import (
	"encoding/json"
	"io"
)

type ChannelMember struct {
	Id        string `json:"string"`
	UserId    string `json:"string"`
	ChannelId string `json:"string"`
}

func (o *ChannelMember) ToJson() string {
	b, err := json.Marshal(o)
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}

func ChannelMemberFromJson(data io.Reader) *ChannelMember {
	decoder := json.NewDecoder(data)
	var o ChannelMember
	err := decoder.Decode(&o)
	if err == nil {
		return &o
	} else {
		return nil
	}
}
