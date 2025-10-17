package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/teamcutter/tempest/internal/sensorpb"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := sensorpb.NewSensorServiceClient(conn)

	for {
		data := &sensorpb.SensorData{
			DeviceId:    fmt.Sprintf("sensor-%03d", rand.Intn(1000)),
			Temperature: 20 + rand.Float64()*15,
			Humidity:    30 + rand.Float64()*50,
			Timestamp:   time.Now().UnixMilli(),
		}

		_, err := client.SendData(context.Background(), data)
		if err != nil {
			fmt.Println("Send error:", err)
		}

		time.Sleep(500 * time.Millisecond)
	}
}
