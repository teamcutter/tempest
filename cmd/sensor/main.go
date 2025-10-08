package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/teamcutter/tempest/internal/model"
	"github.com/teamcutter/tempest/internal/pulsar_client"
)

func main() {
	client := pulsar_client.NewClient()
	defer client.Close()

	producer, _ := client.CreateProducer(pulsar.ProducerOptions{
		Topic: "sensors",
	})
	defer producer.Close()

	for {
		sensorData := model.SensorData{
			DeviceID:    fmt.Sprintf("sensor-%03d", rand.Intn(1000)),
			Temperature: 20 + rand.Float64()*15,
			Humidity:    30 + rand.Float64()*50,
			Timestamp:   time.Now().UnixMilli(),
		}
		payload, _ := json.Marshal(sensorData)

		producer.Send(context.Background(), &pulsar.ProducerMessage{
			Payload: payload,
		})

		time.Sleep(500 * time.Millisecond)
	}
}
