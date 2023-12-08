package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"pop_v1/models"
	"pop_v1/router"
	"pop_v1/utils"
	"time"

	// mydisc "pop_v1/discovery"

	//"pop_v1/router"

	//"github.com/gofiber/fiber/v2"

	"github.com/libp2p/go-libp2p"
	gostream "github.com/libp2p/go-libp2p-gostream"
	p2phttp "github.com/libp2p/go-libp2p-http"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/tecbot/gorocksdb"
)

// DiscoveryInterval is how often we re-publish our mDNS records.
const DiscoveryInterval = time.Hour

// DiscoveryServiceTag is used in our mDNS advertisements to discover other  peers.
const DiscoveryServiceTag = "pubsub"

const TopicName = "BlockMagix"
const TopicClient = "Transaction"

var (
	topicNameFlag = flag.String("topicName", "BlockMagix", "name of topic to join")
)

func main() {
	// parse some flags to set our nickname and the room to join
	var name string
	flag.StringVar(&name, "name", "", " name to distinguish nodes")
	flag.Parse()

	if name == "" {
		log.Fatal("You need to specify '-name' of your node to either 'client' , 'superior' or 'validator'")
	}
	ctx := context.Background()
	// create a new libp2p Host that listens on a random TCP port
	serverHost, err := libp2p.New(libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/0"))
	if err != nil {
		panic(err)
	}

	log.Printf("Host ID: %s", serverHost.ID())
	log.Printf("Connect to me on:")
	for _, addr := range serverHost.Addrs() {
		log.Printf("  %s/p2p/%s", addr, serverHost.ID())
	}
	fmt.Printf("\n")

	// Initialize self host id and multiaddr
	utils.Self_id = serverHost.ID()
	utils.Self_addr = serverHost.Addrs()[0]

	//setup local mDNS discovery
	if err := setupDiscovery(serverHost); err != nil {
		panic(err)
	}
	//go discover.DiscoverPeers(ctx, serverHost, topicNameFlag)
	// create a new PubSub service using the GossipSub router
	ps, err := pubsub.NewGossipSub(ctx, serverHost)
	if err != nil {
		panic(err)
	}
	// join the room from the cli flag, or the flag default
	utils.Topic, err = ps.Join(TopicName)
	if err != nil {
		panic(err)
	}
	utils.CTopic, err = ps.Join(TopicClient)
	if err != nil {
		panic(err)
	}
	// if err := utils.Topic.Publish(context.Background(), []byte("hello")); err != nil {
	// 	fmt.Println("### Publish error:", err)
	// }
	listener, _ := gostream.Listen(serverHost, p2phttp.DefaultP2PProtocol)
	defer listener.Close()
	go func() {
		router.MainRoute()
		server := &http.Server{}
		server.Serve(listener)
	}()

	// // Start server
	// server := &http.Server{}
	// server.Serve(listener)

	sub, err := utils.Topic.Subscribe()

	if err != nil {
		panic(err)
	}

	subC, errC := utils.CTopic.Subscribe()

	if errC != nil {
		panic(err)
	}
	go printMessagesFrom(ctx, sub)
	addTransactions(ctx, subC)
}

func addTransactions(ctx context.Context, sub *pubsub.Subscription) {
	for {
		m, err := sub.Next(ctx)
		if err != nil {
			panic(err)
		}
		var transaction models.TransactionInfo
		fmt.Println("Transaction Received From : ", m.ReceivedFrom)
		err = json.Unmarshal([]byte(m.Message.Data), &transaction)
		if err != nil {
			log.Fatal(err)
		}
		key := utils.GenerateIdClient()
		options := gorocksdb.NewDefaultOptions()
		options.SetCreateIfMissing(true)

		new_transaction := models.NewTransaction(key, transaction.Sender_id, transaction.Receiver_id, transaction.Amount)

		db, err := gorocksdb.OpenDb(options, "database/transaction")
		if err != nil {
			log.Fatal(err)
		}

		// // Serialize the transaction to JSON
		transactionData, err := json.Marshal(new_transaction)
		if err != nil {
			fmt.Printf("Error serializing the transaction: %v\n", err)
			return
		}
		value := []byte(transactionData)

		writeOptions := gorocksdb.NewDefaultWriteOptions()
		defer writeOptions.Destroy()

		err = db.Put(writeOptions, []byte(key), value)
		if err != nil {
			fmt.Printf("Error storing the transaction: %v\n", err)
			return
		}
		fmt.Println("trasaction added Successfully to client")
		db.Close()
		ShowTransactions()
	}
}

func printMessagesFrom(ctx context.Context, sub *pubsub.Subscription) {
	for {
		m, err := sub.Next(ctx)
		if err != nil {
			panic(err)
		}
		fmt.Println(m.ReceivedFrom, ": ", string(m.Message.Data))
	}
}

// app := fiber.New()
// router.MainRoute(app)
// serverAddr := "0.0.0.0:" + config.Config("PORT")
// fmt.Printf("Starting server on %s\n", serverAddr)
// if err := app.Listen(serverAddr); err != nil {
// 	log.Fatalf("Server failed to start: %v\n", err)
// }

//app := fiber.New()
//router.MainRoute(app)

// Start the Fiber app in a goroutine
// go func() {
// 	serverAddr := "0.0.0.0:8888" // Change the port as needed
// 	log.Printf("Fiber app listening on %s\n", serverAddr)
// 	err := app.Listen(serverAddr)
// 	if err != nil {
// 		log.Fatalf("Fiber app failed to start: %v\n", err)
// 	}
// }()
// setup peer discovery
// go mydisc.DiscoverPeers(ctx, h, &config.room)

// discoveryNotifee gets notified when we find a new peer via mDNS discovery
type discoveryNotifee struct {
	h host.Host
}

// HandlePeerFound connects to peers discovered via mDNS. Once they're connected,
// the PubSub system will automatically start interacting with them if they also
// support PubSub.
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	fmt.Printf("discovered new peer %s\n", pi.ID)
	err := n.h.Connect(context.Background(), pi)
	if err != nil {
		fmt.Printf("error connecting to peer %s: %s\n", pi.ID, err)
	}
}

// setupDiscovery creates an mDNS discovery service and attaches it to the libp2p Host.
// This lets us automatically discover peers on the same LAN and connect to them.
func setupDiscovery(h host.Host) error {
	// setup mDNS discovery to find local peers
	s := mdns.NewMdnsService(h, DiscoveryServiceTag, &discoveryNotifee{h: h})
	return s.Start()
}

// package main

// import (
// 	"fmt"
// 	"log"
// 	"pop_v1/config"
// 	"pop_v1/router"

// 	"github.com/gofiber/fiber/v2"
// )

//func main() {

// app := fiber.New()
// router.MainRoute(app)
// serverAddr := "0.0.0.0:"+config.Config("PORT")
// fmt.Printf("Starting server on %s\n", serverAddr)
// if err := app.Listen(serverAddr); err != nil {
// 	log.Fatalf("Server failed to start: %v\n", err)
// }

// }

// func(w http.ResponseWriter, r *http.Request) {
// 	fmt.Printf("%v\n", r.Method)
// 	switch r.Method {
// 	case http.MethodPost:
// 		// Handle POST request
// 		// body, err := ioutil.ReadAll(r.Body)
// 		// if err != nil {
// 		// 	http.Error(w, "Error reading request body", http.StatusInternalServerError)
// 		// 	return
// 		// }

// 		// // Do something with the POST data
// 		// fmt.Printf("Received POST request with body: %s\n", body)

// 		// // Respond to the client
// 		// w.Write([]byte("Received the POST data successfully"))

// 		// // Optionally, publish the data to a topic
// 		// if err := utils.Topic.Publish(context.Background(), body); err != nil {
// 		// 	fmt.Println("### Publish error:", err)
// 		// }
// 		fmt.Print("HELLO I AM CALLING IN POST")
// 		w.Write([]byte("Received the post jjhbjhhjdata successfully"))
// 	case http.MethodGet:
// 		w.Write([]byte("Received the get jjhbjhhjdata successfully"))

// 	default:
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 	}
// })

// =========================================> code for the testing purpose

func ShowTransactions() {
	// Open a RocksDB database
	opts := gorocksdb.NewDefaultOptions()
	defer opts.Destroy()
	opts.SetCreateIfMissing(true)
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
