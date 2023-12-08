package noderouter

import (
	"net/http"
	"pop_v1/utils"
)

// SetupRoutes func
// func SetupHealthRoute(app *fiber.App) {
// 	client := app.Group("/")
// 	client.Head("/alive", utils.Health)
// }

// SetupRoutes func
func SetupHealthRoute() {

	http.HandleFunc("/alive", utils.Health)
}
