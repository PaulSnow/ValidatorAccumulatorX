package database

import (
	"fmt"
	"os"

	"github.com/dgraph-io/badger/v2"

	"github.com/PaulSnow/LoadTest/organizedDataAccumulator/types"
)

type DB struct {
	DBHome   string
	badgerDB *badger.DB
}

func (d *DB) Init(instance int) {
	// Make sure the home directory exists.
	d.DBHome = fmt.Sprintf("%s%s%03d", types.GetHomeDir(), "/.ValAcc/badger", instance)
	os.MkdirAll(d.DBHome, 0777)
	db, err := badger.Open(badger.DefaultOptions(d.DBHome))
	if err != nil {
		panic(err)
	}
	d.badgerDB = db
}

// GetKey
// Given a bucket and a key, return the combined key
func GetKey(bucket string, key []byte) (CKey []byte) {
	CKey = append(CKey, []byte(bucket)...)
	CKey = append(CKey, key...)
	return CKey
}

// Get
// Look in the given bucket, and return the key found.  Returns nil if no value
// is found for the given key
func (d *DB) Get(bucket string, key []byte) (value []byte) {
	CKey := GetKey(bucket, key) // combine the bucket and the key

	// Go look up the CKey, and return any error we might find.
	err := d.badgerDB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(CKey)
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			value = append(value, val...)
			return nil
		})
		return err
	})
	// If anything goes wrong, return nil
	if err != nil {
		return nil
	}
	// If we didn't find the value, we will return a nil here.
	return value
}

// Set
// Set a key/value in the database.  We return an error if there was a problem
// writing the key/value pair to the database.
func (d *DB) Set(bucket string, key []byte, value []byte) error {
	CKey := GetKey(bucket, key)

	// Update the key/value in the database
	err := d.badgerDB.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(CKey), value)
		return err
	})
	return err
}
