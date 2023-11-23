package responsecontroller

import (
	"fmt"
	"net/http"
	"pop_v1/models"
	"pop_v1/utils"
	"pop_v1/config"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

var lock sync.Mutex

func Recieveresponse(c *fiber.Ctx) error {
	var response models.Response
	fmt.Print("HELLOfsd")
	if err := c.BodyParser(&response); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	fmt.Println("recieved response")
	if models.T == 1 {
		return c.Status(408).JSON(fiber.Map{
			"message": "Timeout",
		})
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
	lock.Unlock()
	if models.TotalResponse == 0 {
		client := http.Client{
			Timeout: 5 * time.Millisecond,
		}
		lock.Lock()
	models.TotalResponse++
	lock.Unlock()
		url := "http://"+config.Config("CLIENT")+"/signal"
		client.Head(url)
	}else {
		lock.Lock()
		models.TotalResponse++
		lock.Unlock()

	}
	
	fmt.Printf("responsecount:-->%v\n", models.TotalResponse)
	return c.JSON("response")
}
