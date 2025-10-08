package pulsar_client

import (
	"log"

	"github.com/apache/pulsar-client-go/pulsar"
)

func NewClient() pulsar.Client {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: "pulsar://pulsar:6650",
	})
	if err != nil {
		log.Fatalf("Cannot connect to Pulsar: %v", err)
	}
	return client
}
