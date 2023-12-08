package transactionhelper

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"pop_v1/models"

	"github.com/tecbot/gorocksdb"
)

// func Gettransaction(c *fiber.Ctx) error {
// 	options := gorocksdb.NewDefaultOptions()
// 	options.SetCreateIfMissing(true)

// 	db, err := gorocksdb.OpenDb(options, "database/transaction")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()

// 	readOptions := gorocksdb.NewDefaultReadOptions()
// 	defer readOptions.Destroy()

// 	iter := db.NewIterator(readOptions)
// 	defer iter.Close()

// 	transactions := make([]models.Transaction, 0)

// 	for iter.SeekToFirst(); iter.Valid(); iter.Next() {
// 		//key := iter.Key()
// 		value := iter.Value()
// 		transactionData := value.Data()

// 		// Deserialize transaction from JSON
// 		var transaction models.Transaction
// 		if err := json.Unmarshal(transactionData, &transaction); err != nil {
// 			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
// 		}

// 		transactions = append(transactions, transaction)
// 	}

// 	return c.JSON(transactions)
// }

func Gettransaction(w http.ResponseWriter, r *http.Request) {
	options := gorocksdb.NewDefaultOptions()
	options.SetCreateIfMissing(true)

	db, err := gorocksdb.OpenDb(options, "database/transaction")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	readOptions := gorocksdb.NewDefaultReadOptions()
	defer readOptions.Destroy()

	iter := db.NewIterator(readOptions)
	defer iter.Close()

	transactions := make([]models.Transaction, 0)

	for iter.SeekToFirst(); iter.Valid(); iter.Next() {
		//key := iter.Key()
		value := iter.Value()
		transactionData := value.Data()

		// Deserialize transaction from JSON
		var transaction models.Transaction
		if err := json.Unmarshal(transactionData, &transaction); err != nil {
			http.Error(w, "Transaction retrival error", http.StatusNotFound)
			return
		}

		transactions = append(transactions, transaction)
	}
	for i := range transactions {
		fmt.Println(transactions[i])

	}

}
