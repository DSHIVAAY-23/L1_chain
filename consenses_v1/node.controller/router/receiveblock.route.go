package noderouter

import (
	"net/http"
	responsecontroller "pop_v1/node.controller/response.controller"
)

// SetupRoutes func
//
//	func NodeRoutes(app *fiber.App) {
//		client := app.Group("/")
//		client.Post("/receive", responsecontroller.Receive)
//		client.Post("/groupresponse",responsecontroller.Groupresponse)
//		client.Head("/signalgroup",responsecontroller.Groupsignal)
//	}
func NodeRoutes() {
	http.HandleFunc("/recieve", responsecontroller.Receive)
	http.HandleFunc("/groupresponse", responsecontroller.Groupresponse)
	http.HandleFunc("/signalgroup",responsecontroller.Groupsignal)


}
