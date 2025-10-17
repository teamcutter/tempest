package main

import (
	"context"
	"fmt"
	"net"

	"net/http"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/teamcutter/tempest/internal/pulsar_client"
	"github.com/teamcutter/tempest/internal/sensorpb"
	"google.golang.org/grpc"
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

type SensorServer struct{
	sensorpb.UnimplementedSensorServiceServer
	producer pulsar.Producer
}

func (s *SensorServer) SendData(ctx context.Context, data *sensorpb.SensorData) (*sensorpb.SensorResponse, error) {
	msgCount.Inc()

	if data.Temperature > 30 {
		highTemp.Inc()
		alert := fmt.Sprintf("High temperature alert! Sensor ID: %s, Temp: %.2f", data.DeviceId, data.Temperature)
		s.producer.Send(ctx, &pulsar.ProducerMessage{
			Payload: []byte(alert),
		})
		fmt.Println("Sent alert:", alert) 
	}

	return &sensorpb.SensorResponse{Status: "ok"}, nil
}

func init() {
	prometheus.MustRegister(msgCount)
	prometheus.MustRegister(highTemp)
}

func main() {
	client := pulsar_client.NewClient()
	defer client.Close()

	producer, _ := client.CreateProducer(pulsar.ProducerOptions{
		Topic: "alerts",
	})
	defer producer.Close()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	sensorpb.RegisterSensorServiceServer(server, &SensorServer{producer: producer})
	fmt.Println("gRPC server listening on :50051")
	if err := server.Serve(lis); err != nil {
		panic(err)
	}
}
