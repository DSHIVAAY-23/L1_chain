package transactionhelper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"pop_v1/models"
	"pop_v1/utils"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/tecbot/gorocksdb"
)

func Recievetransaction(c *fiber.Ctx) error {
	// Parse JSON data from the request
	var transaction models.TransactionInfo
	fmt.Println("HELLO")
	if err := c.BodyParser(&transaction); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	options := gorocksdb.NewDefaultOptions()
	options.SetCreateIfMissing(true)

	new_transaction := models.NewTransaction(GenerateRandomString(10), transaction.Sender_id, transaction.Receiver_id, transaction.Amount)

	db, err := gorocksdb.OpenDb(options, "database/transaction")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Serialize the transaction to JSON
	transactionData, err := json.Marshal(new_transaction)
	if err != nil {
		fmt.Printf("Error serializing the transaction: %v\n", err)
		return c.Status(500).JSON("error in serializing the transaction")
	}

	key := []byte(new_transaction.Transaction_id)
	value := []byte(transactionData)

	writeOptions := gorocksdb.NewDefaultWriteOptions()
	defer writeOptions.Destroy()

	err = db.Put(writeOptions, key, value)
	if err != nil {
		fmt.Printf("Error storing the transaction: %v\n", err)
		return c.Status(500).JSON("error in storing the data in db ")
	}
	fmt.Println("trasaction added Successfully to client")
	// code to broadcast

	// Define a list of destination IP addresses (replace with your own IP addresses).
	destIPs, err := utils.GetNodeIps()
	if err != nil {
		fmt.Printf("Error in getting the ips of nodes in network: %v\n", err)
		return c.Status(500).JSON("Error in getting the ips of nodes in network:")
	}

	// Use a WaitGroup to wait for all requests to complete.
	var wg sync.WaitGroup

	// Use a Mutex to safely access shared data.
	var mu sync.Mutex

	// Define a function to send the message to a destination IP.
	sendMessage := func(destIP string) {
		defer wg.Done()

		client := &http.Client{}
		req, err := http.NewRequest("POST", "http://"+destIP+"/recievetransaction", bytes.NewBuffer(transactionData))
		if err != nil {
			fmt.Printf("Error creating request for %s: %v\n", destIP, err)
			return
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			mu.Lock()
			fmt.Printf("Error sending message to %s: %v\n", destIP, err)
			mu.Unlock()
			return
		}
		defer resp.Body.Close()

		mu.Lock()
		fmt.Printf("Message sent to %s, Status Code: %v\n", destIP, resp.Status)
		mu.Unlock()
	}

	// Send messages concurrently to all destinations.
	for _, destIP := range destIPs {
		wg.Add(1)
		go sendMessage(destIP)
	}

	wg.Wait()

	return c.Status(201).JSON(transaction)
}
