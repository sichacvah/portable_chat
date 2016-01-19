package model

import (
	"encoding/json"
	"io"
)

const (
	POST_SYSTEM_MESSAGE_PREFIX = "system_"
	POST_JOIN_LEAVE            = "system_join_leave"
)

type Post struct {
	UserId    string            `json:"user_id"`
	ChannelId string            `json:"channel_id"`
	Id        string            `json:"id"`
	Message   string            `json:"message"`
	Props     map[string]string `json:"props"`
}

func (o *Post) ToJson() string {
	b, err := json.Marshal(o)
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}

func PostFromJson(data io.Reader) *Post {
	decoder := json.NewDecoder(data)
	var o Post
	err := decoder.Decode(&o)
	if err == nil {
		return &o
	} else {
		return nil
	}
}

func (o *Post) IsValid() *AppError {

	if len(o.Id) == 0 {
		return NewAppError("Post.IsValid", "Invalid Id", "")
	}

	if len(o.UserId) == 0 {
		return NewAppError("Post.UserId", "Invalid UserId", "")
	}

	return nil
}
