package noderouter

import (
	responsecontroller "pop_v1/node.controller/response.controller"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes func
func NodeRoutes(app *fiber.App) {
	client := app.Group("/")
	client.Post("/receive", responsecontroller.Receive)
	client.Post("/groupresponse",responsecontroller.Groupresponse)
	client.Head("/signalgroup",responsecontroller.Groupsignal)
}
