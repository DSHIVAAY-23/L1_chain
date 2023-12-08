package transactionhelper

import (
	"encoding/json"
	"fmt"
	"log"
	"pop_v1/models"

	"github.com/tecbot/gorocksdb"
)

func PrintAllTransactions() {
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
    for iter.SeekToFirst(); iter.Valid(); iter.Next() {
        // key := iter.Key()
        value := iter.Value()
        transactionData := value.Data()

        // Deserialize transaction from JSON
        var transaction models.Transaction
        if err := json.Unmarshal(transactionData, &transaction); err != nil {
            log.Printf("Error: %s", err.Error())
        }


        // // Print the details of each transaction to the terminal
        fmt.Printf("Transaction ID: %s\n", transaction.Transaction_id)
        fmt.Printf("Sender ID: %s\n", transaction.Sender_id)
        fmt.Printf("Receiver ID: %s\n", transaction.Receiver_id)
        fmt.Printf("Amount: %d\n", transaction.Amount)
        fmt.Printf("Timestamp: %s\n", transaction.Timestamp)
        fmt.Println("---------------")
		// fmt.Println(transactions)
    }
}