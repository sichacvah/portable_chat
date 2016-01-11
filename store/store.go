package store

import "github.com/sichacvah/portable_chat/model"

type StoreResult struct {
	Data interface{}
	Err  *model.AppError
}

type StoreChannel chan StoreResult

type Store interface {
	User() UserStore
	Close()
}

// type ChannelStore interface {
// 	Save(channel *model.Channel) StoreChannel
// 	Get(id string) StoreChannel
// 	GetByName(name string) StoreChannel
// 	Delete(id string) StoreChannel
// }

type UserStore interface {
	Save(user *model.User) StoreChannel
	Get(id string) StoreChannel
	Update(user *model.User) StoreChannel
	GetByLogin(login string) StoreChannel
	GetUsers() StoreChannel
	Delete(userId string) StoreChannel
}
