package utils

import (
	"encoding/json"
	"fmt"
	"pop_v1/models"

	"github.com/tecbot/gorocksdb"
)

func GetNodeIps() ([]string, error) {
	// Open a RocksDB database
	opts := gorocksdb.NewDefaultOptions()
	defer opts.Destroy()
	db, err := gorocksdb.OpenDb(opts, "database/node-ip")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return nil, err
	}
	defer db.Close()
	// Reading data
	readOpts := gorocksdb.NewDefaultReadOptions()
	defer readOpts.Destroy()
	// Iterating through the database
	iter := db.NewIterator(readOpts)
	defer iter.Close()
	var ips []string
	for iter.SeekToFirst(); iter.Valid(); iter.Next() {
		//	key := iter.Key()
		value := iter.Value()
		client := models.Nodeinfo{}
		if err := json.Unmarshal(value.Data(), &client); err != nil {
			fmt.Println("Error deserializing data", err)
			return nil, err
		}
		ips = append(ips,client.Addr+"/p2p/"+client.ID)

	}
	return ips, nil
}
