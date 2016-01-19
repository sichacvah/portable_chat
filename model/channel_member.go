package model

import (
	"encoding/json"
	"io"
)

const (
	CHANNEL_ROLE_ADMIN = "channel_admin"
	CHANNEL_ROLE_USER  = "channel_user"
)

type ChannelMember struct {
	UserId    string `json:"string"`
	ChannelId string `json:"string"`
	Role      string `json:"string"`
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

func ChannelMembersMapFromJson(data io.Reader) map[*ChannelMember]bool {
	decoder := json.NewDecoder(data)
	o := make(map[*ChannelMember]bool)
	err := decoder.Decode(&o)
	if err == nil {
		return o
	} else {
		return nil
	}
}
