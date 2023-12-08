package responsecontroller

import (
	"fmt"
	"net/http"
	"pop_v1/models"
	finalizecontroller "pop_v1/node.controller/finalize.controller"
	"time"
)

func Groupsignal(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Signal recieved by admin")
	time.Sleep(10 * time.Second)
	fmt.Print("Came!")
	models.Lock.Lock()
	models.T = 1
	models.Lock.Unlock()
	finalizecontroller.Finalize_response()

}
