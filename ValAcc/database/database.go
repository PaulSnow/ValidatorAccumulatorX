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
	//"fmt"
	//"os"

	"github.com/AccumulusNetwork/ValidatorAccumulator/ValAcc/types"
	dbm "github.com/tendermint/tm-db"
)

type DB struct {
	DBHome   string
	db2 dbm.DB
}



// We take an instance of the database, because we anticipate sometime in the future,
// running multiple instances of the database.  This feature might not ever be used
// for the ValAcc project, but it has been useful for factomd testing.
func (d *DB) InitDB(db dbm.DB) {
    d.db2 = db
}

func (d *DB) Init(instance int) {
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
	value, err := d.db2.Get(CKey)
	if err != nil {
		return nil
	}
	return value
}

func (d *DB) GetInt32(bucket string, ikey uint32) (value []byte) {
	key := types.Uint32Bytes(ikey)
	return d.Get(bucket, key)
}

// Put
// Put a key/value in the database.  We return an error if there was a problem
// writing the key/value pair to the database.
func (d *DB) Put(bucket string, key []byte, value []byte) error {
	CKey := GetKey(bucket, key)
	return d.db2.Set(CKey,value)
}

// PutInt
// Put a key/value in the database, where the key is an index.  We return an error if there was a problem
// writing the key/value pair to the database.
func (d *DB) PutInt32(bucket string, ikey int, value []byte) error {
	key := types.Uint32Bytes(uint32(ikey))
	return d.Put(bucket, key, value)
}
