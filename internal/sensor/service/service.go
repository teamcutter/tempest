package service

import (
	"context"
	"fmt"
	"log"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/teamcutter/tempest/internal/sensorpb"
	"github.com/teamcutter/tempest/internal/shared"
)

type SensorServer struct {
	sensorpb.UnimplementedSensorServiceServer
	Producer pulsar.Producer
}

func (s *SensorServer) SendData(ctx context.Context, data *sensorpb.SensorData) (*sensorpb.SensorResponse, error) {
	shared.MsgCount.Inc()

	if data.Temperature > 30 {
		shared.HighTemp.Inc()
		alert := fmt.Sprintf("High temperature alert! Sensor ID: %s, Temp: %.2f", data.DeviceId, data.Temperature)
		s.Producer.Send(ctx, &pulsar.ProducerMessage{
			Payload: []byte(alert),
		})
		log.Println("Sent alert:", alert)
	}

	return &sensorpb.SensorResponse{Status: "ok"}, nil
}
