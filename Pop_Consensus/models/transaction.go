package models

import (
	"strconv"
	"time"
	"github.com/tecbot/gorocksdb"
)

type TransactionInfo struct {
	Sender_id   string
	Receiver_id string
	Amount      int
}
type Transaction struct {
	Transaction_id string
	Sender_id      string
	Receiver_id    string
	Amount         int
	Timestamp      string // You can use a Unix timestamp
}

func NewTransaction(transaction_id, sender_id, receiver_id string, amount int) Transaction {

	return Transaction{
		Transaction_id: transaction_id,
		Sender_id:      sender_id,
		Receiver_id:    receiver_id,
		Amount:         amount,
		Timestamp:      strconv.FormatInt(time.Now().Unix(), 10),
	}
}

func GetLastTransactionID() string {
	options := gorocksdb.NewDefaultOptions()
	options.SetCreateIfMissing(true)

	db, err := gorocksdb.OpenDb(options, "database/transaction")
	if err != nil {
		//return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		return ""
	}
	defer db.Close()

	readOptions := gorocksdb.NewDefaultReadOptions()
	defer readOptions.Destroy()

	// Create a new iterator and seek to the last key
	iterator := db.NewIterator(readOptions)
	defer iterator.Close()

	iterator.SeekToLast()

	if iterator.Valid() {
		// Extract the last transaction ID and return it
		key := iterator.Key()
		return string(key.Data())
	}

	return ""
}


