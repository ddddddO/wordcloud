package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cloud.google.com/go/pubsub"
)

// NOTE: ローカル上のエミュレータに接続できないよう。。
func main() {
	fmt.Println("launch")

	projectID := "test-emu-555"
	subID := "my_subscription"
	if err := pullMsgsSync(os.Stdout, projectID, subID); err != nil {
		fmt.Errorf("%v", err)
	}
}

func pullMsgsSync(w io.Writer, projectID, subID string) error {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	sub := client.Subscription(subID)
	exists, err := sub.Exists(ctx)
	if err != nil {
		return err
	}
	if !exists {
		fmt.Println("non exist subscription")
		return nil
	}

	// Turn on synchronous mode. This makes the subscriber use the Pull RPC rather
	// than the StreamingPull RPC, which is useful for guaranteeing MaxOutstandingMessages,
	// the max number of messages the client will hold in memory at a time.
	sub.ReceiveSettings.Synchronous = true
	sub.ReceiveSettings.MaxOutstandingMessages = 1

	// Receive messages for 5 seconds.
	ctx, cancel := context.WithTimeout(ctx, 180*time.Second)
	defer cancel()

	// Create a channel to handle messages to as they come in.
	cm := make(chan *pubsub.Message)
	defer close(cm)
	// Handle individual messages in a goroutine.
	go func() {
		for msg := range cm {
			fmt.Fprintf(w, "Got message :%q\n", string(msg.Data))
			msg.Ack()
		}
	}()

	fmt.Println("block..")

	// Receive blocks until the passed in context is done.
	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		fmt.Println("received?")
		cm <- msg
	})
	if err != nil && status.Code(err) != codes.Canceled {
		return fmt.Errorf("Receive: %v", err)
	}

	fmt.Println("block....")

	return nil
}
