package store

import (
	"log"

	"github.com/joyrexus/buckets"
)

type BoltDBCli struct {
	db *buckets.DB
}

var instanceBx *BoltDBCli = nil

func Init() *BoltDBCli {
	if instanceBx == nil {
		db, err := buckets.Open("fourty_four.db")
		if err != nil {
			log.Fatal(err)
		}
		instanceBx = &BoltDBCli{db: db}
	}
	return instanceBx
}

func Connect() *BoltDBCli {
	return instanceBx
}

func (boltdb *BoltDBCli) SetValue(bucketName string, key string, value []byte) error {
	bx, err := instanceBx.db.New([]byte(bucketName))
	if err != nil {
		return err
	}
	bx.Put([]byte(key), value)
	return err
}

func (boltdb *BoltDBCli) GetValue(bucketName string, key string) (data []byte, err error) {
	bx, err := instanceBx.db.New([]byte(bucketName))
	if err != nil {
		data = []byte("")
		return
	}
	data, err = bx.Get([]byte(key))
	return
}

func (boltdb *BoltDBCli) DeleteValue(bucketName string, key string) error {
	bx, err := instanceBx.db.New([]byte(bucketName))
	if err != nil {
		return err
	}
	bx.Delete([]byte(key))
	return nil
}

func (boltdb *BoltDBCli) ListValues(bucketName string) ([][]byte, error) {
	listValues := [][]byte{}
	bx, err := instanceBx.db.New([]byte(bucketName))
	if err != nil {
		return listValues, err
	}
	items, err := bx.Items()
	if err != nil {
		return listValues, err
	}

	for _, item := range items {
		listValues = append(listValues, item.Value)
	}
	return listValues, err
}
