package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"pop_v1/models"

	"github.com/tecbot/gorocksdb"
)

// // add ip and port of all nodes to the the db
// func AddNode(c *fiber.Ctx) error {
// 	// get a Id
// 	id := GenerateId()

// 	// Open a RocksDB database
// 	opts := gorocksdb.NewDefaultOptions()
// 	defer opts.Destroy()
// 	opts.SetCreateIfMissing(true)
// 	db, err := gorocksdb.OpenDb(opts, "database/node-ip")
// 	if err != nil {
// 		fmt.Println("Error opening database:", err)
// 		return c.Status(500).JSON(fiber.Map{
// 			"message": "DB Error while opening the database !",
// 		})
// 	}
// 	defer db.Close()

// 	//Fetching the client data from the post body
// 	client := models.Node{}
// 	c.BodyParser(&client)

// 	// Generating the next id and Serielizing the struct
// 	clientJSON, err := json.Marshal(client)
// 	if err != nil {
// 		fmt.Println("Error serializing client:", err)
// 		return c.Status(500).JSON(fiber.Map{
// 			"Error": err,
// 		})
// 	}

// 	// Writing data to the db
// 	writeOpts := gorocksdb.NewDefaultWriteOptions()
// 	defer writeOpts.Destroy()
// 	err = db.Put(writeOpts, []byte(id), clientJSON)
// 	if err != nil {
// 		fmt.Println("Error writing data:", err)
// 		return c.Status(500).JSON(fiber.Map{
// 			"Error": err,
// 		})
// 	}
// 	fmt.Println("node added successfully")
// 	// Success Response
// 	return c.Status(200).JSON(fiber.Map{
// 		"message": "Client added successfully!",
// 	})
// }

// func GetNode(c *fiber.Ctx) error {
// 	// Open a RocksDB database
// 	opts := gorocksdb.NewDefaultOptions()
// 	defer opts.Destroy()
// 	db, err := gorocksdb.OpenDb(opts, "database/node-ip")
// 	if err != nil {
// 		fmt.Println("Error opening database:", err)
// 		return c.Status(500).JSON(fiber.Map{
// 			"message": "DB Error while opening the database!",
// 		})
// 	}
// 	defer db.Close()

// 	// Reading data
// 	readOpts := gorocksdb.NewDefaultReadOptions()
// 	defer readOpts.Destroy()

// 	// Iterating through the database
// 	iter := db.NewIterator(readOpts)
// 	defer iter.Close()

// 	var ips []struct {
// 		ID      string
// 		Address models.Node
// 	}

// 	for iter.SeekToFirst(); iter.Valid(); iter.Next() {
// 		key := iter.Key()
// 		value := iter.Value()
// 		client := models.Node{}
// 		if err := json.Unmarshal(value.Data(), &client); err != nil {
// 			fmt.Println("Error deserializing data", err)
// 			return c.Status(500).JSON(fiber.Map{
// 				"Err": err,
// 			})
// 		}
// 		ips = append(ips, struct {
// 			ID      string
// 			Address models.Node
// 		}{
// 			ID:      string(key.Data()),
// 			Address: client,
// 		})
// 	}

// 	return c.Status(200).JSON(ips)
// }

// add ip and port of all nodes to the the db
func AddNode(w http.ResponseWriter, r *http.Request) {
	// get a Id
	println("hi1")
	id := GenerateId()

	// Open a RocksDB database
	opts := gorocksdb.NewDefaultOptions()
	defer opts.Destroy()
	opts.SetCreateIfMissing(true)
	db, err := gorocksdb.OpenDb(opts, "database/node-ip")
	if err != nil {
		fmt.Println("Error opening database:", err)
		return
	}
	defer db.Close()

	//Fetching the node data from the  body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	println("hi2")

	var node models.Nodeinfo

	err = json.Unmarshal(bodyBytes, &node)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\n", node)

	// Generating the next id and Serielizing the struct
	nodeJSON, err := json.Marshal(node)
	if err != nil {
		fmt.Println("Error serializing node", err)
		return
	}
	println("hi3")

	// Writing data to the db
	writeOpts := gorocksdb.NewDefaultWriteOptions()
	defer writeOpts.Destroy()
	err = db.Put(writeOpts, []byte(id), nodeJSON)
	if err != nil {
		fmt.Println("Error writing data: in no-ip db", err)
		return
	}
	println("hi4")

	fmt.Println("node added successfully")
	// Success Response
	return

}

func GetNode(w http.ResponseWriter, r *http.Request) {
	// Open a RocksDB database
	opts := gorocksdb.NewDefaultOptions()
	defer opts.Destroy()
	db, err := gorocksdb.OpenDb(opts, "database/node-ip")
	if err != nil {
		fmt.Println("Error opening database: node-ip", err)
		return
	}
	defer db.Close()

	// Reading data
	readOpts := gorocksdb.NewDefaultReadOptions()
	defer readOpts.Destroy()

	// Iterating through the database
	iter := db.NewIterator(readOpts)
	defer iter.Close()

	var ips []struct {
		ID      string
		Address models.Nodeinfo
	}

	for iter.SeekToFirst(); iter.Valid(); iter.Next() {
		key := iter.Key()
		value := iter.Value()
		node := models.Nodeinfo{}
		if err := json.Unmarshal(value.Data(), &node); err != nil {
			fmt.Println("Error deserializing data in node-ip", err)
			return
		}
		ips = append(ips, struct {
			ID      string
			Address models.Nodeinfo
		}{
			ID:      string(key.Data()),
			Address: node,
		})
	}
	for i := range ips {
		fmt.Println("id->"+ips[i].ID)
		fmt.Println("Addr->"+ips[i].Address.ID)
		fmt.Println("Addr->"+ips[i].Address.Addr)

	}

}
