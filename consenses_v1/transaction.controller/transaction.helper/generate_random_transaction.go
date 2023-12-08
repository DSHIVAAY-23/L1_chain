package transactionhelper

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"pop_v1/models"
	"strconv"
	"github.com/tecbot/gorocksdb"
)

func Generate_random_transaction() {
	options := gorocksdb.NewDefaultOptions()
	options.SetCreateIfMissing(true)
	db, err := gorocksdb.OpenDb(options, "database/transaction")
	if err != nil {
		log.Fatal(err)
	}

	writeOptions := gorocksdb.NewDefaultWriteOptions()
	defer db.Close()
	defer writeOptions.Destroy()

	for i := 0; i < 100; i++ {
		new_transaction := models.NewTransaction(GenerateRandomString(10),("Amit" + strconv.Itoa(i)), "Ayush"+strconv.Itoa(i), rand.Int()*1000)
		// Serialize the transaction to JSON
		transactionData, err := json.Marshal(new_transaction)
		if err != nil {
			fmt.Printf("Error serializing the transaction: %v\n", err)
			return
		}

		key := []byte(new_transaction.Transaction_id)
		value := []byte(transactionData)

		err = db.Put(writeOptions, key, value)
		if err != nil {
			fmt.Printf("Error storing the transaction: %v\n", err)
			return
		}

		fmt.Println("transaction num", i)

	}

}
