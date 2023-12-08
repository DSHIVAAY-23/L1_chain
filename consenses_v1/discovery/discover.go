package discovery

import (
	"context"
	"fmt"

	mydht "pop_v1/dht"

	host "github.com/libp2p/go-libp2p/core/host"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
)

func DiscoverPeers(ctx context.Context, h host.Host, topicNameFlag *string) {
	kademliaDHT := mydht.InitDHT(ctx, h)
	routingDiscovery := drouting.NewRoutingDiscovery(kademliaDHT)
	dutil.Advertise(ctx, routingDiscovery, *topicNameFlag)

	// Look for others who have announced and attempt to connect to them
	anyConnected := false
	for !anyConnected {
		fmt.Println("Searching for peers...")
		peerChan, err := routingDiscovery.FindPeers(ctx, *topicNameFlag)
		if err != nil {
			panic(err)
		}
		for peer := range peerChan {
			if peer.ID == h.ID() {
				continue // No self connection
			}
			err := h.Connect(ctx, peer)
			if err == nil {
				fmt.Println("Connected to:", peer.ID)
				anyConnected = true
			}
			// } else {
			// 	fmt.Printf("Failed connecting to %s, error: %s\n", peer.ID, err)

			// }
		}
	}
	fmt.Println("Peer discovery complete")
}
