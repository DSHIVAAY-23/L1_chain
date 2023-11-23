package noderouter

import (
	"pop_v1/utils"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes func
func SetupHealthRoute(app *fiber.App) {
	client := app.Group("/")
	client.Head("/alive", utils.Health)
}
