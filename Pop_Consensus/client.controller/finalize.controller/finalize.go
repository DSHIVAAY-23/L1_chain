package finalizecontroller

import (
	"fmt"
	"pop_v1/models"
	"sync"
)

var lock sync.Mutex

func Finalize_block() bool {
	var finalize_hash string = ""
	var max int = 0
	for key, value := range models.VoteMap {
		if len(value) >= max {
			finalize_hash = key
		}

	}
	responses := models.VoteMap[finalize_hash]
	for i := range responses {
		if !responses[i].Vote {
			fmt.Println("false")
			lock.Lock()
			models.T = 0
			models.TotalResponse = 0
			models.VoteMap = make(map[string][]models.Response)
			lock.Unlock()
			return false
		}

	}
	lock.Lock()
	models.T = 0
	models.TotalResponse = 0
	models.VoteMap = make(map[string][]models.Response)
	lock.Unlock()
	fmt.Printf("%v", responses[0])
	fmt.Println("yes")
	fmt.Println(models.T)

	return true
}
