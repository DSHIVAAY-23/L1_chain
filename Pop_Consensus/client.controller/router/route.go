package clientrouter

import (
	responsecontroller "pop_v1/client.controller/response.controller"
	"pop_v1/middleware"
	"pop_v1/utils"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes func
func SetupClientRoute(app *fiber.App) {
	client := app.Group("")
	client.Post("/add", middleware.Nodecheck, utils.AddNode)
	client.Get("/get", utils.GetNode)
	client.Post("/recieveresponse", responsecontroller.Recieveresponse)
	client.Head("/signal", responsecontroller.Sendsignal)
}
