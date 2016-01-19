package store

import (
	"strconv"
	"strings"

	"github.com/joyrexus/buckets"
	"github.com/sichacvah/portable_chat/model"
)

type BoltChannelStore struct {
	channelsBucket       *buckets.Bucket
	channelMembersBucket *buckets.Bucket
}

const (
	CHANNELS        = "channels"
	CHANNEL_MEMBERS = "channel_members"
)

func NewBoltDbChannelStore(boltStore *BoltDBStore) ChannelStore {
	cs := &BoltChannelStore{}
	cs.channelsBucket, _ = boltStore.db.New([]byte(CHANNELS))
	cs.channelMembersBucket, _ = boltStore.db.New([]byte(CHANNEL_MEMBERS))
	return cs
}

func (cs BoltChannelStore) GetCount() StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		var result StoreResult
		items, err := cs.channelsBucket.Items()
		if err != nil {
			result.Err = model.NewAppError("BoltChannelStore.GetCount", "Error while get items", "")
		} else {
			result.Data = len(items)
		}
		storeChannel <- result
		close(storeChannel)
		return
	}()

	return storeChannel
}

func (cs BoltChannelStore) SaveDirectChannel(channel *model.Channel, mb1 *model.ChannelMember, mb2 *model.ChannelMember) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		var result StoreResult

		err := cs.channelsBucket.Put([]byte(channel.Id), []byte(channel.ToJson()))
		if err != nil {
			result.Err = model.NewAppError("BoltChannelStore.SaveDirectChannel", "Error while save channel", "")
		} else {

			items, err := cs.channelMembersBucket.Items()
			if err != nil {
				result.Err = model.NewAppError("BoltChannelStore.SaveDirectChannel", "Error while save channel", "")
			}

			count := len(items)
			err = cs.channelMembersBucket.Put([]byte(strconv.Itoa(count+1)), []byte(mb1.ToJson()))
			if err != nil {
				result.Err = model.NewAppError("BoltChannelStore.SaveDirectChannel", "Error while save channel member", "")
			}
			err = cs.channelMembersBucket.Put([]byte(strconv.Itoa(count+2)), []byte(mb2.ToJson()))
			if err != nil {
				result.Err = model.NewAppError("BoltChannelStore.SaveDirectChannel", "Error while save channel member", "")
			}
		}

		result.Data = channel

		storeChannel <- result
		close(storeChannel)
		return
	}()

	return storeChannel
}

func (cs BoltChannelStore) SaveMember(member *model.ChannelMember) StoreChannel {
	storeChannel := make(StoreChannel)
	go func() {
		var result StoreResult
		items, err := cs.channelMembersBucket.Items()
		if err != nil {
			result.Err = model.NewAppError("BoltChannelStore.SaveMember", "Error while save members", "")
		} else {
			id := len(items) + 1
			err := cs.channelMembersBucket.Put([]byte(strconv.Itoa(id)), []byte(member.ToJson()))
			if err != nil {
				result.Err = model.NewAppError("BoltChannelStore.SaveMember", "Error while save members", "")
			} else {
				result.Data = member
			}
		}

		storeChannel <- result
		close(storeChannel)
		return
	}()
	return storeChannel
}

func (cs BoltChannelStore) GetMember(channelId string, userId string) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		var result StoreResult
		var member *model.ChannelMember
		var itemString string
		notFound := true

		items, err := cs.channelMembersBucket.Items()
		if err != nil {
			result.Err = model.NewAppError("BoltChannelStore.GetMember", "Error while get members", "")
		} else {
			for _, item := range items {
				itemString = string(item.Value)
				member = model.ChannelMemberFromJson(strings.NewReader(itemString))
				if member.ChannelId == channelId && member.UserId == userId {
					result.Data = member
					notFound = false
				}
			}
			if notFound {
				result.Err = model.NewAppError("BoltChannelStore.GetMember", "Not found", "")
			}
		}

		storeChannel <- result
		close(storeChannel)
		return
	}()

	return storeChannel
}

func (cs BoltChannelStore) GetMembers(channel *model.Channel) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		var result StoreResult

		items, err := cs.channelMembersBucket.Items()
		if err != nil {
			result.Err = model.NewAppError("BoltChannelStore.GetMembers", "Error while get members", "")
		} else {
			resultData := make(map[*model.ChannelMember]bool)
			for _, item := range items {
				memberJson := string(item.Value)
				member := model.ChannelMemberFromJson(strings.NewReader(memberJson))
				if member.ChannelId == channel.Id {
					resultData[member] = true
				}
			}

			result.Data = resultData
		}

		storeChannel <- result
		close(storeChannel)
		return
	}()

	return storeChannel
}

func (cs BoltChannelStore) Save(channel *model.Channel) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		var result StoreResult
		if channel.Type == model.CHANNEL_DIRECT {
			result.Err = model.NewAppError("BoltChannelStore.Save", "Use Direct channel save to save direct channel", "")
		} else {
			channel.PreSave()
			err := cs.channelsBucket.Put([]byte(channel.Id), []byte(channel.ToJson()))
			if err != nil {
				result.Err = model.NewAppError("BoltChannelStore.Save", "Error while save", "")
			}
		}

		storeChannel <- result
		close(storeChannel)
		return
	}()

	return storeChannel
}

func (cs BoltChannelStore) Update(channel *model.Channel) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		var result StoreResult
		channelJson := channel.ToJson()
		err := cs.channelsBucket.Put([]byte(channel.Id), []byte(channelJson))
		if err != nil {
			result.Err = model.NewAppError("BoltChannelStore.Update", "Error while update", "")
		} else {
			result.Data = channelJson
		}
		storeChannel <- result
		close(storeChannel)
		return
	}()

	return storeChannel
}

func (cs BoltChannelStore) Delete(channelId string) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		var result StoreResult

		if len(channelId) <= 0 {
			result.Err = model.NewAppError("BoltChanelStore.Delete", "You must get channelId in delete", "")
			storeChannel <- result
			close(storeChannel)
			return
		}
		err := cs.channelsBucket.Delete([]byte(channelId))
		if err != nil {
			result.Err = model.NewAppError("BoltUserStore.Delete", err.Error(), "")
			storeChannel <- result
			close(storeChannel)
		}
	}()

	return storeChannel
}

func (cs BoltChannelStore) GetByName(name string) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		var result StoreResult

		items, err := cs.channelsBucket.Items()
		if err != nil {
			result.Err = model.NewAppError("BoltChannelStore.GetByName", "Error while get by name", "")
			storeChannel <- result
			close(storeChannel)
			return
		}

		for _, item := range items {
			channel := model.ChannelFromJson(strings.NewReader(string(item.Value)))
			if channel.Name == name {
				result.Data = channel
				storeChannel <- result
				close(storeChannel)
				return
			}
		}
	}()

	return storeChannel
}

func (cs BoltChannelStore) Get(channelId string) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		var result StoreResult
		channelJson, err := cs.channelsBucket.Get([]byte(channelId))
		if err != nil {
			result.Err = model.NewAppError("BoltChannelStore.Get", "Error while get", "")
		} else {
			result.Data = model.ChannelFromJson(strings.NewReader(string(channelJson)))
		}

		storeChannel <- result
		close(storeChannel)
		return
	}()

	return storeChannel
}

func (cs BoltChannelStore) GetChannelMembers(channelId string) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		var result StoreResult
		var channelMember *model.ChannelMember
		items, err := cs.channelsBucket.Items()

		if err != nil {
			result.Err = model.NewAppError("BoltChannelStore.GetChannelMembers", "Error while get members", "")
		} else {
			data := make(map[*model.ChannelMember]bool)
			for _, item := range items {
				itemString := string(item.Value)
				channelMember = model.ChannelMemberFromJson(strings.NewReader(itemString))
				if channelMember != nil && channelId == channelMember.ChannelId {
					data[channelMember] = true
				}
			}
			result.Data = data
		}
		storeChannel <- result
		close(storeChannel)
		return
	}()

	return storeChannel
}

func (cs BoltChannelStore) GetChannels(userId string) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		var result StoreResult
		var channel *model.Channel

		items, err := cs.channelsBucket.Items()
		if err != nil {
			result.Err = model.NewAppError("BoltChannelStore.GetChannels", "Error while get items", "")
		} else {
			data := make(map[*model.Channel]bool)

			for _, item := range items {
				channel = model.ChannelFromJson(strings.NewReader(string(item.Value)))
				if channel != nil {
					data[channel] = true
				}
			}
		}

	}()

	return storeChannel
}
