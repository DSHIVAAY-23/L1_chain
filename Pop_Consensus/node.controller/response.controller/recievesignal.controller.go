package responsecontroller

import (
	"fmt"
	"pop_v1/models"
	finalizecontroller "pop_v1/node.controller/finalize.controller"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Groupsignal(c *fiber.Ctx) error {
	fmt.Println("Signal recieved by admin")
	time.Sleep(10 * time.Second)
	fmt.Print("Came!")
	models.Lock.Lock()
	models.T = 1
	models.Lock.Unlock()
	finalizecontroller.Finalize_response()
	return c.Status(200).JSON("SUCCESS")
}
