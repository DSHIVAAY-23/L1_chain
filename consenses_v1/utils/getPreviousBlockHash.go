package utils

import (
	"fmt"

	"github.com/tecbot/gorocksdb"
)

func PreviousHash() string {
	//// Open a RocksDB database
	opts := gorocksdb.NewDefaultOptions()
	defer opts.Destroy()
	opts.SetCreateIfMissing(true)
	db, err := gorocksdb.OpenDb(opts, "database/blockchain")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return "0x000000"
	}
	defer db.Close()

	//// Reading data
	readOpts := gorocksdb.NewDefaultReadOptions()
	defer readOpts.Destroy()
	iter := db.NewIterator(readOpts)
	defer iter.Close()
	iter.SeekToLast()
	if iter.Valid() {
		key := iter.Key()
		return string(key.Data())
	} else {
		return "0xabcdef0000"
	}
}
