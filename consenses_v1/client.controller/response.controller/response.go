package responsecontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"pop_v1/models"
	"pop_v1/utils"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p"
	p2phttp "github.com/libp2p/go-libp2p-http"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

var lock sync.Mutex

// func Recieveresponse(c *fiber.Ctx) error {
// 	var response models.Response
// 	fmt.Print("HELLOfsd")
// 	if err := c.BodyParser(&response); err != nil {
// 		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
// 	}
// 	fmt.Println("recieved response")
// 	if models.T == 1 {
// 		return c.Status(408).JSON(fiber.Map{
// 			"message": "Timeout",
// 		})
// 	}
// 	fmt.Println("YES_RECEIVED")
// 	lock.Lock()
// 	block_recieved := response.Block
// 	block_hash := utils.GenerateBlockHash(block_recieved)
// 	response_list := models.VoteMap[block_hash]
// 	if response_list != nil {
// 		response_list = append(response_list, response)
// 		models.VoteMap[block_hash] = response_list

// 	} else {
// 		var responses []models.Response
// 		responses = append(responses, response)
// 		models.VoteMap[block_hash] = responses
// 	}
// 	lock.Unlock()
// 	if models.TotalResponse == 0 {
// 		client := http.Client{
// 			Timeout: 5 * time.Millisecond,
// 		}
// 		lock.Lock()
// 	models.TotalResponse++
// 	lock.Unlock()
// 		url := "http://"+config.Config("CLIENT")+"/signal"
// 		client.Head(url)
// 	}else {
// 		lock.Lock()
// 		models.TotalResponse++
// 		lock.Unlock()

// 	}

// 	fmt.Printf("responsecount:-->%v\n", models.TotalResponse)
// 	return c.JSON("response")
// }

func Recieveresponse(w http.ResponseWriter, r *http.Request) {
	var response models.Response
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("recieved response by validator")
	if models.T == 1 {
		http.Error(w, "message : Timeout", http.StatusNotFound)
		return
	}
	fmt.Println("YES_RECEIVED")
	lock.Lock()
	block_recieved := response.Block
	block_hash := utils.GenerateBlockHash(block_recieved)
	response_list := models.VoteMap[block_hash]
	if response_list != nil {
		response_list = append(response_list, response)
		models.VoteMap[block_hash] = response_list

	} else {
		var responses []models.Response
		responses = append(responses, response)
		models.VoteMap[block_hash] = responses
	}
	totalResponses := models.TotalResponse
	models.TotalResponse++
	lock.Unlock()
	if totalResponses == 0 {
		clientHost, _ := libp2p.New(libp2p.NoListenAddrs)
		addr := utils.Self_addr
		peerId := utils.Self_id
		info := peer.AddrInfo{
			ID:    peerId,
			Addrs: []multiaddr.Multiaddr{addr},
		}
		clientHost.Connect(context.Background(), info)
		tr := &http.Transport{}
		tr.RegisterProtocol("libp2p", p2phttp.NewTransport(clientHost))

		// specified timeout for the client
		client := &http.Client{
			Transport: tr,
			Timeout:   5 * time.Millisecond,
		}
		lock.Lock()
		fmt.Printf("response-->%v\n", models.TotalResponse)
		lock.Unlock()
		go func() {
			_, err := client.Head("libp2p://" + info.ID.String() + "/signal")
			if err != nil {
				log.Fatalf(err.Error())
			}
		}()

		fmt.Printf("responsecount:-->%v\n", models.TotalResponse)

	} else {
		lock.Lock()
		lock.Unlock()
		fmt.Printf("responsecount:-->%v\n", models.TotalResponse)

	}

}
