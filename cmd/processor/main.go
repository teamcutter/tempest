package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/teamcutter/tempest/internal/model"
	"github.com/teamcutter/tempest/internal/pulsar_client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	msgCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "tempest_messages_total",
		Help: "Amount of processed messages",
	})
	highTemp = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "tempest_high_temp_total",
		Help: "Amount of messages with high temperature",
	})
)

func init() {
	prometheus.MustRegister(msgCount)
	prometheus.MustRegister(highTemp)
}

func main() {
	client := pulsar_client.NewClient()
	defer client.Close()

	consumer, _ := client.Subscribe(pulsar.ConsumerOptions{
		Topic: "sensors",
		SubscriptionName: "processor-sub",
		Type:  pulsar.Shared,
	})
	defer consumer.Close()

	producer, _ := client.CreateProducer(pulsar.ProducerOptions{
		Topic: "alerts",
	})
	defer producer.Close()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	for {
		msg, err := consumer.Receive(context.Background())
		if err != nil {
			fmt.Println("Error receiving message:", err)
			continue
		}

		var sensorData model.SensorData
		if err = json.Unmarshal(msg.Payload(), &sensorData); err == nil {
			msgCount.Inc()
			if sensorData.Temperature > 30.0 {
				highTemp.Inc()
				alert := fmt.Sprintf("High temperature alert! Sensor ID: %s, Temperature: %.2f", sensorData.DeviceID, sensorData.Temperature)
				producer.Send(context.Background(), &pulsar.ProducerMessage{
					Payload: []byte(alert),
				})
				fmt.Println("Sent alert:", alert)
			}
		}
		consumer.Ack(msg)

		time.Sleep(500 * time.Millisecond)
	}
}
