package finalizecontroller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"pop_v1/config"
	"pop_v1/models"
	"time"
)

func Finalize_response() {
	client_ip := config.Config("CLIENT")
	fmt.Print("GO1")
	var max_true models.Block
	var max_false models.Block
	total := float64(models.VoteCount)
	total = total * 0.70
	var max float64
	for _, value := range models.Groupmap {
		if value.Vote {
			max_true = value.Block
			max++
		} else {
			max_false = value.Block
		}
	}
	fmt.Print("GO2")
	var response models.Response
	if max >= total {

		response.Block = max_true
		response.Ip = config.Config("HOST") + ":" + config.Config("PORT")
		response.Vote = true
	} else {
		var response models.Response
		response.Block = max_false
		response.Vote = false
	}
	response.Ip = config.Config("HOST") + ":" + config.Config("PORT")
	fmt.Print("GO3")
	requestBody, err := json.Marshal(response)
	fmt.Printf("\n%+v\n", response)
	if err != nil {
		println("Error in serializing the request")
	}
	client := &http.Client{
		Timeout: time.Second * 2,
	}
	url := "http://" + client_ip + "/recieveresponse"
	fmt.Printf("\n%v\n", url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody)) // Use nil for request body or set requestBody
	fmt.Print("GO4\n")
	if err != nil {
		fmt.Println("message -> Error in Creating request to the client:")
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	models.Lock.Lock()
	models.T=0
	models.Groupmap=make(map[string]models.Response)
	models.VoteCount=0

	models.Lock.Unlock()
	if err != nil {
	fmt.Printf("Error %v with response %v\n", err, resp)
	}
	defer resp.Body.Close()
}
