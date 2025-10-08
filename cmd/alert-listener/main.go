package main

import (
	"context"
	"fmt"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/teamcutter/tempest/internal/pulsar_client"
)

func main() {
	client :=  pulsar_client.NewClient()
	defer client.Close()

	consumer, _ := client.Subscribe(pulsar.ConsumerOptions{
		Topic:            "alerts",
		SubscriptionName: "alert-sub",
		Type:             pulsar.Shared,
	})
	defer consumer.Close()

	fmt.Println("Alert listener started...")
	for {
		msg, err := consumer.Receive(context.Background())
		if err == nil {
			fmt.Println(string(msg.Payload()))
			consumer.Ack(msg)
		}
	}
}