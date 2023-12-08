package responsecontroller

import (
	"fmt"
	"net/http"
	finalizecontroller "pop_v1/client.controller/finalize.controller"
	"pop_v1/models"
	"time"
)

// func Sendsignal(c *fiber.Ctx) error {
// 	fmt.Println("Signal recieved")
// 	time.Sleep(10 * time.Second)
// 	lock.Lock()
// 	models.T = 1
// 	lock.Unlock()
// 	finalizecontroller.Finalize_block()
// 	return c.Status(200).JSON("SUCCESS")
// }

func Sendsignal(http.ResponseWriter, *http.Request) {
	fmt.Println("Signal recieved")
	time.Sleep(10 * time.Second)
	lock.Lock()
	models.T = 1
	lock.Unlock()
	go finalizecontroller.Finalize_block()
	println("Successfully recieved signal")

}
