package finalizecontroller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"pop_v1/config"
	"pop_v1/models"
	"pop_v1/utils"
	"time"

	"github.com/libp2p/go-libp2p"
	p2phttp "github.com/libp2p/go-libp2p-http"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

func Finalize_response() {
	client_id := config.Config("CLIENTID")
	client_addr := config.Config("CLIENTADDR")
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
		response.Id = utils.Self_id.String()
		response.Vote = true
	} else {
		var response models.Response
		response.Block = max_false
		response.Vote = false
	}
	clientHost, _ := libp2p.New(libp2p.NoListenAddrs)
	// orginal admins id and addr
	addr, _ := multiaddr.NewMultiaddr(client_addr)
	peerId, _ := peer.Decode(client_id)
	info := peer.AddrInfo{
		ID:    peerId,
		Addrs: []multiaddr.Multiaddr{addr},
	}
	clientHost.Connect(context.Background(), info)
	tr := &http.Transport{}
	tr.RegisterProtocol("libp2p", p2phttp.NewTransport(clientHost))

	client := &http.Client{
		Transport: tr,
		Timeout:   20 * time.Millisecond,
	}

	response.Id = utils.Self_id.String()
	fmt.Print("GO3")
	requestBody, err := json.Marshal(response)
	fmt.Printf("\n%+v\n", response)
	if err != nil {
		println("Error in serializing the request")
	}
	_, err = client.Post("libp2p://"+info.ID.String()+"/recieveresponse", "application/json", bytes.NewReader(requestBody))
	if err != nil {
		fmt.Println("message -> Error in Creating request to the client:")
		return
	}
	models.Lock.Lock()
	models.T = 0
	models.Groupmap = make(map[string]models.Response)
	models.VoteCount = 0

	models.Lock.Unlock()
	
	fmt.Println("Response send to the client Successfully by admin")
}
