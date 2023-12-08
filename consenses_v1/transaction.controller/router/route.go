package transactionrouter

import (
	"net/http"
	transactionhelper "pop_v1/transaction.controller/transaction.helper"
)

// SetupRoutes func
// func TransactionRoute(app *fiber.App) {
// 	client := app.Group("")
// 	//client.Post("/recievetransaction", transactionhelper.Recievetransaction)
// 	client.Get("/gettransaction", transactionhelper.Gettransaction)
// 	client.Get("/taketransactions",transactionhelper.TransactionLogic)

// }

func TransactionRoute() {
	http.HandleFunc("/recievetransaction", transactionhelper.Recievetransaction)
	http.HandleFunc("/gettransaction",transactionhelper.Gettransaction)
	// function related to superior yet to be changed
	http.HandleFunc("/taketransaction",transactionhelper.TransactionLogic)
}
