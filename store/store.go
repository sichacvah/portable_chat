package store

import "github.com/sichacvah/portable_chat/model"

type StoreResult struct {
	Data interface{}
	Err  *model.AppError
}

type StoreChannel chan StoreResult

type Store interface {
	User() UserStore
	Channel() ChannelStore
	Close()
}

type ChannelStore interface {
	Save(channel *model.Channel) StoreChannel
	GetMembers(channel *model.Channel) StoreChannel
	GetChannels(userId string) StoreChannel
	Get(id string) StoreChannel
	GetByName(name string) StoreChannel
	Delete(id string) StoreChannel
	SaveMember(member *model.ChannelMember) StoreChannel
	GetChannelMembers(channelId string) StoreChannel
	GetMember(channelId string, userId string) StoreChannel
	SaveDirectChannel(channel *model.Channel, mb1 *model.ChannelMember, mb2 *model.ChannelMember) StoreChannel
	GetCount() StoreChannel
}

type UserStore interface {
	Save(user *model.User) StoreChannel
	Get(id string) StoreChannel
	Update(user *model.User) StoreChannel
	GetByLogin(login string) StoreChannel
	GetUsers() StoreChannel
	Delete(userId string) StoreChannel
	GetCount() StoreChannel
}
