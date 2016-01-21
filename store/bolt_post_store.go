package store

import (
	"strings"

	"github.com/joyrexus/buckets"
	"github.com/sichacvah/portable_chat/model"
)

type BoltPostStore struct {
	postStore      *buckets.Bucket
	postsByChannel *buckets.Bucket
}

const (
	POSTS             = "posts"
	POSTS_BY_CHANNELS = "posts_by_channels"
)

func NewBoltDbPostStore(boltStore *BoltDBStore) PostStore {
	ps := &BoltPostStore{}
	ps.postStore, _ = boltStore.db.New([]byte(POSTS))
	ps.postsByChannel, _ = boltStore.db.New([]byte(POSTS_BY_CHANNELS))

	return ps
}

func (ps BoltPostStore) GetPosts(channelID string) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		var result StoreResult
		var posts map[string]*model.Post
		channelPostsJson, err := ps.postsByChannel.Get([]byte(channelID))
		if err != nil {
			result.Err = model.NewAppError("BoltPostStore.", "Post is not valid", "")
		} else {
			channelPosts := model.BoolMapFromJson(strings.NewReader(string(channelPostsJson)))
			postItems, err := ps.postStore.Items()
			if err != nil {
				result.Err = model.NewAppError("BoltPostStore.", "Post is not valid", "")
				storeChannel <- result
				close(storeChannel)
				return
			}

			for _, postItem := range postItems {
				if channelPosts[string(postItem.Key)] {
					postString := string(postItem.Value)
					posts[string(postItem.Key)] = model.PostFromJson(strings.NewReader(postString))
				}
			}
			result.Data = posts
		}

		storeChannel <- result
		close(storeChannel)
		return
	}()

	return storeChannel
}

func (ps BoltPostStore) Get(postId string) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		var result StoreResult
		post, err := ps.postStore.Get([]byte(postId))
		if err != nil {
			result.Err = model.NewAppError("BoltPostStore.", "Post is not valid", "")
		} else {
			result.Data = model.PostFromJson(strings.NewReader(string(post)))
		}

		storeChannel <- result
		close(storeChannel)
		return
	}()

	return storeChannel
}

func (ps BoltPostStore) Update(post *model.Post) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		var result StoreResult
		postErr := post.IsValid()
		if postErr != nil {
			err := ps.postStore.Put([]byte(post.Id), []byte(post.ToJson()))

			if err != nil {
				result.Err = model.NewAppError("BoltPostStore.", "Post is not valid", "")
			} else {
				result.Data = post
			}
		} else {
			result.Err = postErr
		}

		storeChannel <- result
		close(storeChannel)
		return
	}()

	return storeChannel
}

func (ps BoltPostStore) Save(post *model.Post) StoreChannel {
	storeChannel := make(StoreChannel)

	go func() {
		var result StoreResult
		postErr := post.IsValid()
		if postErr != nil {
			post.PreSave()
			err := ps.postStore.Put([]byte(post.Id), []byte(post.ToJson()))

			postsByChannelsJson, err := ps.postsByChannel.Get([]byte(post.ChannelId))
			postsByChannelsString := string(postsByChannelsJson)
			postsByChannels := model.BoolMapFromJson(strings.NewReader(postsByChannelsString))
			postsByChannels[post.Id] = true

			postsByChannelsString = model.BoolMapToJson(postsByChannels)

			err = ps.postStore.Put([]byte(post.ChannelId), []byte(postsByChannelsString))

			if err != nil {
				result.Err = model.NewAppError("BoltPostStore.", "Post is valid", "")
			} else {
				result.Data = post
			}
		} else {
			result.Err = postErr
		}

		storeChannel <- result
		close(storeChannel)
		return
	}()

	return storeChannel
}
