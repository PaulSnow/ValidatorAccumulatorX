package database

// As a key value store, we are not as yet using too many of the features of any database.
// We use hashes to order entries, which are also organized into buckets, much as LevelDB
// allows.
//
//To use this DB interface, you must allocate a DB
// Then call DB.Init(int)
//
// To set a value in the database, call DB.Put(bucket string, key []byte, value []byte) error
//
// To get a value from the database, call DB.Get(bucket string, key []byte) (value []byte)_
//
// see ValAcc/types/types.go for the constants for bucket names

import (
	"fmt"
	"os"

	"github.com/PaulSnow/ValidatorAccumulator/ValAcc/types"

	"github.com/dgraph-io/badger/v2"
)

type DB struct {
	DBHome   string
	badgerDB *badger.DB
}

// We take an instance of the database, because we anticipate sometime in the future,
// running multiple instances of the database.  This feature might not ever be used
// for the ValAcc project, but it has been useful for factomd testing.
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

// Put
// Put a key/value in the database.  We return an error if there was a problem
// writing the key/value pair to the database.
func (d *DB) Put(bucket string, key []byte, value []byte) error {
	CKey := GetKey(bucket, key)

	// Update the key/value in the database
	err := d.badgerDB.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(CKey), value)
		return err
	})
	return err
}

// PutInt
// Put a key/value in the database, where the key is an index.  We return an error if there was a problem
// writing the key/value pair to the database.
func (d *DB) PutInt(bucket string, ikey int, value []byte) error {
	key := types.Uint32Bytes(uint32(ikey))
	CKey := GetKey(bucket, key)

	// Update the key/value in the database
	err := d.badgerDB.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(CKey), value)
		return err
	})
	return err
}
