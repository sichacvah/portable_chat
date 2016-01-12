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
			}
		}
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
