package transactionhelper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"pop_v1/models"
	"pop_v1/utils"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/tecbot/gorocksdb"
)

func Recievetransaction(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	// var transaction models.TransactionInfo

	// err = json.Unmarshal(bodyBytes, &transaction)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(transaction)

	// =============================================================
	// key := utils.GenerateIdClient()
	// options := gorocksdb.NewDefaultOptions()
	// options.SetCreateIfMissing(true)

	// new_transaction := models.NewTransaction(key, transaction.Sender_id, transaction.Receiver_id, transaction.Amount)

	// db, err := gorocksdb.OpenDb(options, "database/transaction")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Serialize the transaction to JSON
	// transactionData, err := json.Marshal(new_transaction)
	// if err != nil {
	// 	fmt.Printf("Error serializing the transaction: %v\n", err)
	// 	return
	// }
	// value := []byte(transactionData)

	// writeOptions := gorocksdb.NewDefaultWriteOptions()
	// defer writeOptions.Destroy()

	// err = db.Put(writeOptions, []byte(key), value)
	// if err != nil {
	// 	fmt.Printf("Error storing the transaction: %v\n", err)
	// 	return
	// }
	// fmt.Println("trasaction added Successfully to client")
	// db.Close()
	// ShowTransactions()
	go doPublish(context.Background(), utils.CTopic, bodyBytes)
	// go utils.StreamConsoleTo(context.Background(), utils.Topic, transactionData)
	// if err := utils.Topic.Publish(context.Background(), []byte("hello")); err != nil {
	// 	fmt.Println("### Publish error:", err)
	// }
}

func doPublish(ctx context.Context, topic *pubsub.Topic, data []byte) {
	if err := topic.Publish(ctx, data); err != nil {
		fmt.Println("### Publish error:", err)
	}
}

// ================================================>> this function is just for the demo purpose
// ================================================>>
// ================================================>>

func ShowTransactions() {
	// Open a RocksDB database
	opts := gorocksdb.NewDefaultOptions()
	defer opts.Destroy()
	db, err := gorocksdb.OpenDb(opts, "database/transaction")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer db.Close()

	// Reading data
	readOpts := gorocksdb.NewDefaultReadOptions()
	defer readOpts.Destroy()

	// Iterating through the database
	iter := db.NewIterator(readOpts)
	defer iter.Close()

	for iter.SeekToFirst(); iter.Valid(); iter.Next() {
		key := iter.Key()
		value := iter.Value()
		client := models.Transaction{}
		if err := json.Unmarshal(value.Data(), &client); err != nil {
			fmt.Println("Error deserializing data", err)
		}
		fmt.Printf("%v : %v\n", string(key.Data()), client)
	}
}

// }
// func Recievetransaction(c *fiber.Ctx) error {
// 	// Parse JSON data from the request
// 	var transaction models.TransactionInfo
// 	fmt.Println("HELLO")
// 	if err := c.BodyParser(&transaction); err != nil {
// 		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
// 	}

// 	options := gorocksdb.NewDefaultOptions()
// 	options.SetCreateIfMissing(true)

// 	new_transaction := models.NewTransaction(GenerateRandomString(10), transaction.Sender_id, transaction.Receiver_id, transaction.Amount)

// 	db, err := gorocksdb.OpenDb(options, "database/transaction")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()

// 	// Serialize the transaction to JSON
// 	transactionData, err := json.Marshal(new_transaction)
// 	if err != nil {
// 		fmt.Printf("Error serializing the transaction: %v\n", err)
// 		return c.Status(500).JSON("error in serializing the transaction")
// 	}

// 	key := []byte(new_transaction.Transaction_id)
// 	value := []byte(transactionData)

// 	writeOptions := gorocksdb.NewDefaultWriteOptions()
// 	defer writeOptions.Destroy()

// 	err = db.Put(writeOptions, key, value)
// 	if err != nil {
// 		fmt.Printf("Error storing the transaction: %v\n", err)
// 		return c.Status(500).JSON("error in storing the data in db ")
// 	}
// 	fmt.Println("trasaction added Successfully to client")
// 	data, err := json.Marshal(new_transaction)
// 	if err != nil {
// 		fmt.Printf("Error marshaling the transaction: %v\n", err)
// 		return c.Status(500).JSON("Error marshaling the transaction")
// 	}
// 	if err := utils.Topic.Publish(context.Background(), []byte(data)); err != nil {
// 		fmt.Println("### Publish error:", err)
// 	}
// 	// code to broadcast

// 	// Define a list of destination IP addresses (replace with your own IP addresses).
// 	// destIPs, err := utils.GetNodeIps()
// 	// if err != nil {
// 	// 	fmt.Printf("Error in getting the ips of nodes in network: %v\n", err)
// 	// 	return c.Status(500).JSON("Error in getting the ips of nodes in network:")
// 	// }

// 	// // Use a WaitGroup to wait for all requests to complete.
// 	// var wg sync.WaitGroup

// 	// Use a Mutex to safely access shared data.
// 	// var mu sync.Mutex

// 	// // Define a function to send the message to a destination IP.
// 	// sendMessage := func(destIP string) {
// 	// 	defer wg.Done()

// 	// 	client := &http.Client{}
// 	// 	req, err := http.NewRequest("POST", "http://"+destIP+"/recievetransaction", bytes.NewBuffer(transactionData))
// 	// 	if err != nil {
// 	// 		fmt.Printf("Error creating request for %s: %v\n", destIP, err)
// 	// 		return
// 	// 	}

// 	// 	req.Header.Set("Content-Type", "application/json")

// 	// 	resp, err := client.Do(req)
// 	// 	if err != nil {
// 	// 		mu.Lock()
// 	// 		fmt.Printf("Error sending message to %s: %v\n", destIP, err)
// 	// 		mu.Unlock()
// 	// 		return
// 	// 	}
// 	// 	defer resp.Body.Close()

// 	// 	mu.Lock()
// 	// 	fmt.Printf("Message sent to %s, Status Code: %v\n", destIP, resp.Status)
// 	// 	mu.Unlock()
// 	// }

// 	// // Send messages concurrently to all destinations.
// 	// for _, destIP := range destIPs {
// 	// 	wg.Add(1)
// 	// 	go sendMessage(destIP)
// 	// }

// 	// wg.Wait()

// 	return c.Status(201).JSON(new_transaction)
// }

// reader := bufio.NewReader(os.Stdin)
// 	for {
// 		s, err := reader.ReadString('\n')
// 		if err != nil {
// 			panic(err)
// 		}
// 		if err := topic.Publish(ctx, []byte(s)); err != nil {
// 			fmt.Println("### Publish error:", err)
// 		}
// 	}
