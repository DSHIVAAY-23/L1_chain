package clientrouter

import (
	"net/http"
	responsecontroller "pop_v1/client.controller/response.controller"
	"pop_v1/utils"
)

// SetupRoutes func
// func SetupClientRoute(app *fiber.App) {
// 	client := app.Group("")
// 	client.Post("/add", middleware.Nodecheck, utils.AddNode)
// 	client.Get("/get", utils.GetNode)
// 	client.Post("/recieveresponse", responsecontroller.Recieveresponse)
// 	client.Head("/signal", responsecontroller.Sendsignal)
// 	client.Get("/", func(c *fiber.Ctx) error {
// 		return c.SendString("Welcome to the libp2p server!")
// 	})
// }

func SetupClientRoute() {
	http.HandleFunc("/recieveresponse", responsecontroller.Recieveresponse)
	http.HandleFunc("/add", utils.AddNode)
	http.HandleFunc("/get", utils.GetNode)
	http.HandleFunc("/signal", responsecontroller.Sendsignal)
}
