package store

import (
	l4g "github.com/alecthomas/log4go"

	"github.com/joyrexus/buckets"
)

type BoltDBStore struct {
	db      *buckets.DB
	user    UserStore
	channel ChannelStore
	post    PostStore
}

func NewBoltDBStore() *BoltDBStore {
	boltDbStore := BoltDBStore{}

	if db, err := buckets.Open("fourty_four.db"); err != nil {
		l4g.Critical("Not db connection")
		panic(err)
	} else {
		boltDbStore.db = db
		boltDbStore.user = NewBoltDbUserStore(&boltDbStore)
		boltDbStore.channel = NewBoltDbChannelStore(&boltDbStore)
		boltDbStore.post = NewBoltDbPostStore(&boltDbStore)
	}
	return &boltDbStore
}

func (bs BoltDBStore) Channel() ChannelStore {
	return bs.channel
}

func (bs BoltDBStore) User() UserStore {
	return bs.user
}

func (bs BoltDBStore) Post() PostStore {
	return bs.post
}

func (bs BoltDBStore) Close() {
	l4g.Info("Closing BoltStore")
	bs.db.Close()
}
