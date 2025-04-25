package kafka

import (
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaWriters struct {
	UserStatsWriter   *kafka.Writer
	StreamStatsWriter *kafka.Writer
}

func NewKafkaWriters() *KafkaWriters {
	return &KafkaWriters{
		UserStatsWriter: kafka.NewWriter(kafka.WriterConfig{
			Brokers:      []string{"localhost:9092"},
			Topic:        "user-stats",
			Balancer:     &kafka.LeastBytes{},
			Async:        true,
			BatchTimeout: 10 * time.Millisecond,
			WriteTimeout: 50 * time.Millisecond,
			RequiredAcks: int(kafka.RequireOne),
		}),
		StreamStatsWriter: kafka.NewWriter(kafka.WriterConfig{
			Brokers:      []string{"localhost:9092"},
			Topic:        "stream-stats",
			Balancer:     &kafka.LeastBytes{},
			Async:        true,
			BatchTimeout: 10 * time.Millisecond,
			WriteTimeout: 50 * time.Millisecond,
			RequiredAcks: int(kafka.RequireOne),
		}),
	}
}

func (w *KafkaWriters) Close() {
	w.UserStatsWriter.Close()
	w.StreamStatsWriter.Close()
}

// func (k *KafkaWriters) GetWriter() *kafka.Writer {
// 	return k.StreamStatsWriter
// }
