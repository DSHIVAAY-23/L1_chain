package transactionhelper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"pop_v1/config"
	"pop_v1/models"
	superiorhelper "pop_v1/superior.controller/helper"
	"pop_v1/utils"
	"strconv"
	"time"

	"github.com/libp2p/go-libp2p"
	p2phttp "github.com/libp2p/go-libp2p-http"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

// func TransactionLogic(c *fiber.Ctx) error {
// 	//// take the transactions from the memepool
// 	out, done, wg := utils.GetMempoolTxns()
// 	var txnsMetaData []string
// 	var txns []models.Transaction
// 	err := -1
// 	println(err)
// 	go func() {
// 		for {
// 			select {
// 			case res := <-done:
// 				if !res {
// 					//error is there
// 					err = 1
// 				} else {
// 					err = 0
// 				}
// 				return
// 			case fun := <-out:
// 				metadata, txn := fun()
// 				txns = append(txns, txn)
// 				txnsMetaData = append(txnsMetaData, metadata)
// 				(*wg).Done()
// 			}
// 		}
// 	}()

// 	// take ips after sending a ping pong message
// 	out_, done_ := utils.TakeIps()
// 	var ips []string
// 	err_ := -1
// 	go func() {
// 		for {
// 			select {
// 			case res := <-done_:
// 				if !res {
// 					//error is there
// 					err_ = 1
// 				} else {
// 					err_ = 0
// 				}
// 				return
// 			case ip := <-out_:
// 				ips = append(ips, ip)
// 			}
// 		}
// 	}()

// 	//// make sure ips and txns are collected.
// 	for {
// 		if err > -1 {
// 			break
// 		}
// 	}
// 	for {
// 		if err_ > -1 {
// 			break
// 		}
// 	}
// 	if err == 1 || err_ == 1 {
// 		return c.Status(500).JSON(fiber.Map{
// 			"message": "Internal Server Error",
// 		})
// 	}

// 	//// Generate Groups

// 	superiorhelper.GroupGenerator(&ips)
// 	groups, _ := strconv.Atoi(config.Config("GROUPS"))
// 	param, _ := strconv.Atoi(config.Config("PARAM"))

// 	//// select admins
// 	var admins []int
// 	for i := 0; i < param; i++ {
// 		random := rand.Intn(100000000) % groups
// 		admins = append(admins, random)
// 	}

// 	//// create the block
// 	block := utils.NewBlock(txns, txnsMetaData)
// 	hash := utils.GenerateBlockHash(block)
// 	requestBody, Err := json.Marshal(block)
// 	if Err != nil {
// 		fmt.Println("Error marshaling JSON:", err)
// 		return c.Status(500).JSON(fiber.Map{
// 			"message": "Internal Server Error",
// 		})
// 	}

// 	//// brodcast the block
// 	for i := 0; i < groups; i++ {
// 		for j := 0; j < param; j++ {
// 			//i-th group and j-th param
// 			//for the jth parameter admins[j] is the admin
// 			go func(I int, J int) {
// 				ip := ips[I*param+J]
// 				url := "http://" + ip + "/receive"
// 				client := &http.Client{
// 					Timeout: 60 * time.Millisecond,
// 				}
// 				headers := map[string]string{
// 					"Content-Type": "application/json",
// 					"group-id":     strconv.Itoa(I),
// 					"param-id":     strconv.Itoa(J),
// 					"admin-ip":     ips[admins[J]*param+J],
// 					"block-hash":   hash,
// 				}
// 				req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody)) // Use nil for request body or set requestBody
// 				if err != nil {
// 					fmt.Println("Error creating the request:", err)
// 					return
// 				}

// 				// Set multiple headers
// 				for key, value := range headers {
// 					req.Header.Set(key, value)
// 				}

// 				// Send the request
// 				resp, err := client.Do(req)
// 				if err != nil {
// 					fmt.Println("Error sending the request:", err)
// 					return
// 				}
// 				defer resp.Body.Close()

// 				// Check the response status code
// 				if resp.StatusCode == http.StatusOK {
// 					fmt.Println("Request was successful.")
// 				} else {
// 					fmt.Printf("Request failed with status code: %d\n", resp.StatusCode)
// 				}
// 			}(i, j)
// 		}
// 	}

// 	// return the response to the client node
// 	return c.Status(200).JSON(block)
// }

func TransactionLogic(w http.ResponseWriter, r *http.Request) {
	//// take the transactions from the memepool
	println("hello1")
	out, done, wg := utils.GetMempoolTxns()
	var txnsMetaData []string
	var txns []models.Transaction
	err := -1

	go func() {
		for {
			select {
			case res := <-done:
				if !res {
					//error is there
					err = 1
				} else {
					err = 0
				}
				return
			case fun := <-out:
				metadata, txn := fun()
				txns = append(txns, txn)
				txnsMetaData = append(txnsMetaData, metadata)
				(*wg).Done()
			}
		}
	}()
	println("hello2")
	// take ips after sending a ping pong message

	// Take ip me jinda wala code dalna hai abhi
	out_, done_ := utils.TakeIps()
	var ips []models.Nodeinfo
	err_ := -1
	go func() {
		for {
			select {
			case res := <-done_:
				if !res {
					//error is there
					err_ = 1
				} else {
					err_ = 0
				}
				return
			case ip := <-out_:
				ips = append(ips, ip)
			}
		}
	}()

	//// make sure ips and txns are collected.
	for {
		if err > -1 {
			break
		}
	}
	for {
		if err_ > -1 {
			break
		}
	}
	if err == 1 || err_ == 1 {
		println("message : Internal Server Error")
	}

	//// Generate Groups
	println("hello3")
	superiorhelper.GroupGenerator(&ips)
	groups, _ := strconv.Atoi(config.Config("GROUPS"))
	param, _ := strconv.Atoi(config.Config("PARAM"))

	//// select admins
	var admins []int
	for i := 0; i < param; i++ {
		random := rand.Intn(100000000) % groups
		admins = append(admins, random)
	}
	println("hello4")
	//// create the block
	block := utils.NewBlock(txns, txnsMetaData)
	hash := utils.GenerateBlockHash(block)
	// requestBody, Err := json.Marshal(block)
	// if Err != nil {
	// 	fmt.Println("Error marshaling JSON:", err)
	// }

	//// brodcast the block
	for i := 0; i < groups; i++ {
		for j := 0; j < param; j++ {
			//i-th group and j-th param
			//for the jth parameter admins[j] is the admin
			go func(I int, J int) {
				ip := ips[I*param+J]
				clientHost, _ := libp2p.New(libp2p.NoListenAddrs)
				peer_id, _ := peer.Decode(ip.ID)
				peer_addr, _ := multiaddr.NewMultiaddr(ip.Addr)
				info := peer.AddrInfo{
					ID:    peer_id,
					Addrs: []multiaddr.Multiaddr{peer_addr},
				}

				clientHost.Connect(context.Background(), info)
				tr := &http.Transport{}
				tr.RegisterProtocol("libp2p", p2phttp.NewTransport(clientHost))
				client := &http.Client{
					Transport: tr,
					Timeout:   60 * time.Millisecond,
				}
				var request models.Superior_res
				request.Aid = ips[admins[J]*param+J].ID
				request.A_addr = ips[admins[J]*param+J].Addr
				request.Gid = strconv.Itoa(I)
				request.Pid = strconv.Itoa(J)
				request.Block_hash = hash
				request.Block = block

				requestBody, Err := json.Marshal(request)
				if Err != nil {
					fmt.Println("Error marshaling Data: request", err)
				}
				println(info.ID.String())
				println(info.Addrs[0].String())
				_, err := client.Post("libp2p://"+info.ID.String()+"/recieve", "application/json", bytes.NewReader(requestBody))
				println("hello5")
				if err != nil {
					fmt.Printf("Error creating the request: to "+info.ID.String(), err)
					return
				}

			}(i, j)
		}
	}

}
