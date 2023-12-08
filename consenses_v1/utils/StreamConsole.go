package utils

import (
	"context"
	"fmt"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

func StreamConsoleTo(ctx context.Context, topic *pubsub.Topic, data []byte) {

	if err := topic.Publish(ctx, data); err != nil {
		fmt.Println("### Publish error:", err)
	}
}
