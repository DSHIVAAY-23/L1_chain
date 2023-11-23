package responsecontroller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	finalizecontroller "pop_v1/node.controller/finalize.controller"
	"pop_v1/config"
	"pop_v1/models"
	"pop_v1/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Receive(c *fiber.Ctx) error {
	if models.T == 1 {
		return c.Status(500).JSON(fiber.Map{
			"message": "Not accepting response:",
		})
	}
	selfaddr := config.Config("HOST") + ":" + config.Config("PORT")
	gid := c.Get("group-id")
	pid := c.Get("param-id")
	aip := c.Get("admin-ip")
	block_hash := c.Get("block-hash")
	var block models.Block
	if err := c.BodyParser(&block); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	recieved_block_hash := utils.GenerateBlockHash(block)

	var response models.Response
	response.Block = block
	response.Ip = selfaddr

	if recieved_block_hash != block_hash {
		response.Vote = false
	} else {
		if finalizecontroller.CheckParameter(pid, block) {
			response.Vote = true
		} else {
			response.Vote = false
		}

	}

	fmt.Printf("\n\nparam-id->%v\nResponse->%v\n",pid,response)
	
	client := &http.Client{
		Timeout: 20 * time.Millisecond,
	}

	requestBody, err := json.Marshal(response)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "resposne can not be serilaized",
		})
	}
	// node itselif is the admin

	if aip == selfaddr {
		models.Lock.Lock()
		models.Groupmap[selfaddr] = response
		models.VoteCount++
		models.Lock.Unlock()
		client_mini := http.Client{
			Timeout: 5 * time.Millisecond,
		}
		url := "http://" + selfaddr + "/signalgroup"
		client_mini.Head(url)
		fmt.Printf("Group-id : %v\nParam-id : %v\nAdmin-IP : %v\nBlock- : %v\n", gid, pid, aip, block)
		return c.Status(200).JSON(fiber.Map{
			"message": "Success",
		})
	}
	req, err := http.NewRequest("POST", "http://"+aip+"/groupresponse", bytes.NewBuffer(requestBody)) // Use nil for request body or set requestBody
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Error in Creating:",
		})
	}
	req.Header.Set("Content-Type", "application/json")
	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Error sending the request:",
		})
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Request was successful.")
	} else {
		fmt.Printf("Request failed with status code: %d\n", resp.StatusCode)
	}
	fmt.Printf("Group-id : %v\nParam-id : %v\nAdmin-IP : %v\nBlock- : %v\n", gid, pid, aip, block)
	return c.Status(200).JSON(fiber.Map{
		"message": "Success",
	})
}
