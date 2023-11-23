package transactionhelper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"pop_v1/config"
	"pop_v1/models"
	"pop_v1/utils"
	"strconv"
	"time"

	superiorhelper "pop_v1/superior.controller/helper"

	"github.com/gofiber/fiber/v2"
)

func TransactionLogic(c *fiber.Ctx) error {
	//// take the transactions from the memepool
	out, done, wg := utils.GetMempoolTxns()
	var txnsMetaData []string
	var txns []models.Transaction
	err := -1
	go func() {
		for {
			select {
			case res := <-done:
				if !res{
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

	//// take ips after sending a ping pong message
	out_, done_ := utils.TakeIps()
	var ips []string
	err_ := -1
	go func() {
		for {
			select {
			case res := <-done_:
				if !res{
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
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	//// Generate Groups
	
	superiorhelper.GroupGenerator(&ips)
	groups, _ := strconv.Atoi(config.Config("GROUPS"))
	param, _ := strconv.Atoi(config.Config("PARAM"))

	//// select admins
	var admins []int
	for i := 0; i < param; i++ {
		random := rand.Intn(100000000) % groups
		admins = append(admins, random)
	}

	//// create the block
	block := utils.NewBlock(txns, txnsMetaData)
	hash := utils.GenerateBlockHash(block)
	requestBody, Err := json.Marshal(block)
	if Err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return c.Status(500).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	//// brodcast the block
	for i := 0; i < groups; i++ {
		for j := 0; j < param; j++ {
			//i-th group and j-th param
			//for the jth parameter admins[j] is the admin
			go func(I int, J int) {
				ip := ips[I*param+J]
				url := "http://" + ip + "/receive"
				client := &http.Client{
					Timeout: 60* time.Millisecond,
				}
				headers := map[string]string{
					"Content-Type": "application/json",
					"group-id":     strconv.Itoa(I),
					"param-id":     strconv.Itoa(J),
					"admin-ip":     ips[admins[J]*param+J],
					"block-hash":   hash,
				}
				req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody)) // Use nil for request body or set requestBody
				if err != nil {
					fmt.Println("Error creating the request:", err)
					return
				}

				// Set multiple headers
				for key, value := range headers {
					req.Header.Set(key, value)
				}

				// Send the request
				resp, err := client.Do(req)
				if err != nil {
					fmt.Println("Error sending the request:", err)
					return
				}
				defer resp.Body.Close()

				// Check the response status code
				if resp.StatusCode == http.StatusOK {
					fmt.Println("Request was successful.")
				} else {
					fmt.Printf("Request failed with status code: %d\n", resp.StatusCode)
				}
			}(i, j)
		}
	}

	//// return the response to the client node
	return c.Status(200).JSON(block)
}
