package main

import (
	"fmt"
	"net"

	"net/http"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/teamcutter/tempest/internal/pulsar_client"
	"github.com/teamcutter/tempest/internal/sensorpb"
	"github.com/teamcutter/tempest/internal/sensor/service"
	"github.com/teamcutter/tempest/internal/shared"
	"google.golang.org/grpc"
)

func init() {
	prometheus.MustRegister(shared.MsgCount)
	prometheus.MustRegister(shared.HighTemp)
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
	sensorpb.RegisterSensorServiceServer(server, &service.SensorServer{Producer: producer})
	fmt.Println("gRPC server listening on :50051")
	if err := server.Serve(lis); err != nil {
		panic(err)
	}
}
