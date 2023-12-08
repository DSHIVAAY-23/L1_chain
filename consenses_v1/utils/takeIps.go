package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"pop_v1/config"
	"pop_v1/models"
	"strconv"

	"github.com/libp2p/go-libp2p"
	p2phttp "github.com/libp2p/go-libp2p-http"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/tecbot/gorocksdb"
)

func TakeIps() (<-chan models.Nodeinfo, <-chan bool) {
	out := make(chan models.Nodeinfo)
	done := make(chan bool)
	go func() {
		//// Open a RocksDB database
		opts := gorocksdb.NewDefaultOptions()
		defer opts.Destroy()
		opts.SetCreateIfMissing(true)
		db, err := gorocksdb.OpenDb(opts, "database/node-ip")
		if err != nil {
			fmt.Println("Error opening database:", err)
			done <- false
			close(done)
			close(out)
			return
		}
		defer db.Close()
		//// Reading data
		readOpts := gorocksdb.NewDefaultReadOptions()
		defer readOpts.Destroy()
		iter := db.NewIterator(readOpts)
		defer iter.Close()
		cnt := 0
		param, _ := strconv.Atoi(config.Config("PARAM"))
		gps, _ := strconv.Atoi(config.Config("GROUPS"))
		CNT := param * gps
		for iter.SeekToFirst(); iter.Valid(); iter.Next() {
			value := iter.Value()
			node := models.Nodeinfo{}
			if err := json.Unmarshal(value.Data(), &node); err != nil {
				fmt.Println("Error deserializing data", err)
				done <- false
				close(done)
				close(out)
				return
			}
			clientHost, _ := libp2p.New(libp2p.NoListenAddrs)
			peer_id, _ := peer.Decode(node.ID)
			peer_addr, _ := multiaddr.NewMultiaddr(node.Addr)
			info := peer.AddrInfo{
				ID:    peer_id,
				Addrs: []multiaddr.Multiaddr{peer_addr},
			}
			clientHost.Connect(context.Background(), info)
			tr := &http.Transport{}
			tr.RegisterProtocol("libp2p", p2phttp.NewTransport(clientHost))

			//	ip := node.Addr.String() + ":" + fmt.Sprint(node.PORT)

			//	url := "http://" + ip + "/alive"
			// client := &http.Client{
			// 	Transport: tr,
			// 	Timeout:   10 * time.Millisecond,
			// }
			//	res, err := client.Head("libp2p://" + info.ID.String() + "/alive")
			// if err != nil {
			// 	log.Fatalf(err.Error())
			// }
			// res, err := client.Head(url)
			// if err != nil || res.StatusCode != 200 {
			// 	continue
			// }
			out <- node
			cnt++
			if cnt == CNT {
				break
			}
		}
		done <- true
		close(done)
		close(out)
	}()
	return out, done
}
