package store

import (
	l4g "code.google.com/p/log4go"

	"github.com/joyrexus/buckets"
)

type BoltDBStore struct {
	db   *buckets.DB
	user UserStore
}

func NewBoltDBStore() BoltDBStore {
err:
	boltDbStore := &BoltDBStore{}

	db, err := buckets.Open("fourty_four.db"); err != nil {
		l4g.Critical("Not db connection")
		panic(err)
	} else {
		boltDbStore.db = db
		boltDbStore.user = NewBoltDbUser(boltDbStore)
	}
	return boltDbStore
}

func (bs BoltDBStore) User() UserStore {
	return bs.user
}

func (bs BoltDBStore) Close() {
	l4g.Info("Closing BoltStore")
	bs.db.Close()
}
