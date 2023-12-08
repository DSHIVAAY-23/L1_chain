package responsecontroller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"pop_v1/models"
	finalizecontroller "pop_v1/node.controller/finalize.controller"
	"pop_v1/utils"
	"time"

	"github.com/libp2p/go-libp2p"
	p2phttp "github.com/libp2p/go-libp2p-http"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

// func Receive(c *fiber.Ctx) error {
// 	if models.T == 1 {
// 		return c.Status(500).JSON(fiber.Map{
// 			"message": "Not accepting response:",
// 		})
// 	}
// 	selfaddr := config.Config("HOST") + ":" + config.Config("PORT")
// 	gid := c.Get("group-id")
// 	pid := c.Get("param-id")
// 	aip := c.Get("admin-ip")
// 	block_hash := c.Get("block-hash")
// 	var block models.Block
// 	if err := c.BodyParser(&block); err != nil {
// 		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
// 	}
// 	recieved_block_hash := utils.GenerateBlockHash(block)

// 	var response models.Response
// 	response.Block = block
// 	response.Ip = selfaddr

// 	if recieved_block_hash != block_hash {
// 		response.Vote = false
// 	} else {
// 		if finalizecontroller.CheckParameter(pid, block) {
// 			response.Vote = true
// 		} else {
// 			response.Vote = false
// 		}

// 	}

// 	fmt.Printf("\n\nparam-id->%v\nResponse->%v\n",pid,response)

// 	client := &http.Client{
// 		Timeout: 20 * time.Millisecond,
// 	}

// 	requestBody, err := json.Marshal(response)
// 	if err != nil {
// 		return c.Status(500).JSON(fiber.Map{
// 			"message": "resposne can not be serilaized",
// 		})
// 	}
// 	// node itselif is the admin

// 	if aip == selfaddr {
// 		models.Lock.Lock()
// 		models.Groupmap[selfaddr] = response
// 		models.VoteCount++
// 		models.Lock.Unlock()
// 		client_mini := http.Client{
// 			Timeout: 5 * time.Millisecond,
// 		}
// 		url := "http://" + selfaddr + "/signalgroup"
// 		client_mini.Head(url)
// 		fmt.Printf("Group-id : %v\nParam-id : %v\nAdmin-IP : %v\nBlock- : %v\n", gid, pid, aip, block)
// 		return c.Status(200).JSON(fiber.Map{
// 			"message": "Success",
// 		})
// 	}
// 	req, err := http.NewRequest("POST", "http://"+aip+"/groupresponse", bytes.NewBuffer(requestBody)) // Use nil for request body or set requestBody
// 	if err != nil {
// 		return c.Status(500).JSON(fiber.Map{
// 			"message": "Error in Creating:",
// 		})
// 	}
// 	req.Header.Set("Content-Type", "application/json")
// 	// Send the request
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return c.Status(500).JSON(fiber.Map{
// 			"message": "Error sending the request:",
// 		})
// 	}
// 	defer resp.Body.Close()

// 	// Check the response status code
// 	if resp.StatusCode == http.StatusOK {
// 		fmt.Println("Request was successful.")
// 	} else {
// 		fmt.Printf("Request failed with status code: %d\n", resp.StatusCode)
// 	}
// 	fmt.Printf("Group-id : %v\nParam-id : %v\nAdmin-IP : %v\nBlock- : %v\n", gid, pid, aip, block)
// 	return c.Status(200).JSON(fiber.Map{
// 		"message": "Success",
// 	})
// }

func Receive(w http.ResponseWriter, r *http.Request) {
	println("receive")
	if models.T == 1 {
		println("message : Not accepting response:")
		return
	}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var request models.Superior_res

	err = json.Unmarshal(bodyBytes, &request)
	if err != nil {
		log.Fatal(err)
	}

	var block models.Block
	block = request.Block
	recieved_block_hash := utils.GenerateBlockHash(block)

	var response models.Response
	response.Block = block
	response.Id = utils.Self_id.String()

	if recieved_block_hash != request.Block_hash {
		response.Vote = false
	} else {
		if finalizecontroller.CheckParameter(request.Pid, block) {
			response.Vote = true
		} else {
			response.Vote = false
		}

	}

	fmt.Printf("\n\nparam-id->%v\nResponse->%v\n", request.Pid, response)
	clientHost, _ := libp2p.New(libp2p.NoListenAddrs)
	// orginal admins id and addr
	addr, _ := multiaddr.NewMultiaddr(request.A_addr)
	peerId, _ := peer.Decode(request.Aid)
	info := peer.AddrInfo{
		ID:    peerId,
		Addrs: []multiaddr.Multiaddr{addr},
	}
	clientHost.Connect(context.Background(), info)
	tr := &http.Transport{}
	tr.RegisterProtocol("libp2p", p2phttp.NewTransport(clientHost))

	client := &http.Client{
		Transport: tr,
		Timeout:   20 * time.Millisecond,
	}

	requestBody, err := json.Marshal(response)
	if err != nil {
		println("message : resposne can not be serilaized")
		return
	}
	// node itselif is the admin

	if request.Aid == utils.Self_id.String() {
		models.Lock.Lock()
		models.Groupmap[utils.Self_id.String()] = response
		models.VoteCount++
		models.Lock.Unlock()
		clientHostmini, _ := libp2p.New(libp2p.NoListenAddrs)
		addr := utils.Self_addr
		peerId := utils.Self_id
		info := peer.AddrInfo{
			ID:    peerId,
			Addrs: []multiaddr.Multiaddr{addr},
		}
		clientHostmini.Connect(context.Background(), info)
		tr := &http.Transport{}
		tr.RegisterProtocol("libp2p", p2phttp.NewTransport(clientHostmini))

		client_mini := http.Client{
			Transport: tr,
			Timeout:   5 * time.Millisecond,
		}

		_, err := client_mini.Head("libp2p://" + info.ID.String() + "/signalgroup")
		if err != nil {

			println("message : Error in Creating:")
		}

		fmt.Println("Successfully Signaled " + info.ID.String())
		return
	}
	_, err = client.Post("libp2p://"+info.ID.String()+"/groupresponse", "application/json", bytes.NewReader(requestBody))
	if err != nil {
		println("message : Error in Creating:")
	}
	println("message : Successfully send response")
}
