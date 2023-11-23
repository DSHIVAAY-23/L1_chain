package transactionrouter

import (
	
	transactionhelper "pop_v1/transaction.controller/transaction.helper"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes func
func TransactionRoute(app *fiber.App) {
	client := app.Group("")
	client.Post("/recievetransaction", transactionhelper.Recievetransaction)
	client.Get("/gettransaction", transactionhelper.Gettransaction)
	client.Get("/taketransactions",transactionhelper.TransactionLogic)

}