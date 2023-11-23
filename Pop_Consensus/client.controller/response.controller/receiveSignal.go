package responsecontroller

import (
	"fmt"
	"pop_v1/client.controller/finalize.controller"
	"pop_v1/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Sendsignal(c *fiber.Ctx) error {
	fmt.Println("Signal recieved")
	time.Sleep(10 * time.Second)
	lock.Lock()
	models.T = 1
	lock.Unlock()
	finalizecontroller.Finalize_block()
	return c.Status(200).JSON("SUCCESS")
}
