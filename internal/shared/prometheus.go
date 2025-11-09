package shared

import "github.com/prometheus/client_golang/prometheus"

var (
	MsgCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "tempest_messages_total",
		Help: "Amount of processed messages",
	})
	HighTemp = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "tempest_high_temp_total",
		Help: "Amount of messages with high temperature",
	})
)
