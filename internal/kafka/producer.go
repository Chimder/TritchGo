package kafka

import (
	"time"

	"github.com/segmentio/kafka-go"
)

type ProdKafka struct {
	writer *kafka.Writer
}

func NewProdKafka() *ProdKafka {
	return &ProdKafka{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:      []string{"localhost:9092"},
			Topic:        "tritch-stats",
			Balancer:     &kafka.LeastBytes{},
			Async:        true,
			BatchTimeout: 10 * time.Millisecond,
			WriteTimeout: 50 * time.Millisecond,
			RequiredAcks: int(kafka.RequireOne),
			ReadTimeout: 30 * time.Second,
		}),
	}
}

func (k *ProdKafka) GetWriter() *kafka.Writer {
	return k.writer
}

func (k *ProdKafka) Close() {
	k.writer.Close()
}
