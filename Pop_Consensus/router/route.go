package router

import (
	clientrouter "pop_v1/client.controller/router"
	noderouter "pop_v1/node.controller/router"
	transactionrouter "pop_v1/transaction.controller/router"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes func
func MainRoute(app *fiber.App) {
	transactionrouter.TransactionRoute(app)
	clientrouter.SetupClientRoute(app)
	noderouter.NodeRoutes(app)
	noderouter.SetupHealthRoute(app)
	

}