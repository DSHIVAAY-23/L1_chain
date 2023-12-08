package transactionhelper

import (
	"encoding/json"
	"log"
	"pop_v1/models"
	"github.com/tecbot/gorocksdb"
)

func Retrievetransaction(transaction_ids []string) ([]models.Transaction, error) {
	
	var transaction_list []models.Transaction
	
	options := gorocksdb.NewDefaultOptions()
	options.SetCreateIfMissing(true)

	db, err := gorocksdb.OpenDb(options, "database/transaction")
	if err != nil {
		//return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		return nil, err
	}
	defer db.Close()

	readOptions := gorocksdb.NewDefaultReadOptions()
	defer readOptions.Destroy()

	for i := range transaction_ids {
		
		data, err := db.Get(readOptions, []byte(transaction_ids[i]))
		if err != nil {
			log.Fatal("transactions can not be retrieved")
			return nil, err // An error occurred while checking the key
		}

		var transaction models.Transaction
        if err := json.Unmarshal(data.Data(), &transaction); err != nil {
			log.Fatal("transactions can not be retrieved")
            log.Printf("Error: %s", err.Error())
        }

		transaction_list=append(transaction_list, transaction)
		
	}
	return transaction_list,nil
}
