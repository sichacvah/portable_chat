package model

import (
	"encoding/json"
	"io"

	"github.com/pborman/uuid"
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
	CreateAt  int64             `json:"create_at"`
	UpdateAt  int64             `json:"update_at"`
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

func PostsMapToJson(o map[string]*Post) string {
	b, err := json.Marshal(o)
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}

func (o *Post) PreSave() {
	if o.Id == "" {
		o.Id = uuid.New()
	}

	if o.CreateAt == 0 {
		o.CreateAt = GetMillis()
	}

	o.UpdateAt = o.CreateAt

	if o.Props == nil {
		o.Props = make(map[string]string)
	}
}

func (o *Post) IsValid() *AppError {

	if len(o.Id) == 0 {
		return NewAppError("Post.IsValid", "Invalid Id", "")
	}

	if len(o.UserId) == 0 {
		return NewAppError("Post.UserId", "Invalid UserId", "")
	}

	if len(o.ChannelId) == 0 {
		return NewAppError("Post.ChannelID", "Invalid ChannelID", "")
	}

	return nil
}
