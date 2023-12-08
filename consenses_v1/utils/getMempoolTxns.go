package utils

import (
	"encoding/json"
	"fmt"
	"pop_v1/models"
	"sync"

	"github.com/tecbot/gorocksdb"
)

func GetMempoolTxns() (<-chan func() (string, models.Transaction), <-chan bool, *sync.WaitGroup) {
	out := make(chan func() (string, models.Transaction))
	done := make(chan bool)
	wg := sync.WaitGroup{}
	go func() {
		// Open a RocksDB database
		opts := gorocksdb.NewDefaultOptions()
		defer opts.Destroy()
		opts.SetCreateIfMissing(true)
		db, err := gorocksdb.OpenDb(opts, "database/transaction")
		if err != nil {
			fmt.Println("Error opening database:", err)
			done <- false
			close(done)
			close(out)
			return
		}
		defer db.Close()
		// Reading data
		readOpts := gorocksdb.NewDefaultReadOptions()
		defer readOpts.Destroy()
		iter := db.NewIterator(readOpts)
		defer iter.Close()
		cnt := 0
		for iter.SeekToFirst(); iter.Valid(); iter.Next() {
			key := iter.Key()
			value := iter.Value()
			transaction := models.Transaction{}
			if err := json.Unmarshal(value.Data(), &transaction); err != nil {
				fmt.Println("Error deserializing data", err)
				done <- false
				close(done)
				close(out)
				return
			}
			wg.Add(1)
			out <- (func() (string, models.Transaction) {
				return string(key.Data()), transaction
			})
			wg.Wait()
			cnt++
			if cnt == 20 {
				break
			}
		}
		done <- true
		close(done)
		close(out)
	}()
	return out, done, &wg
}
