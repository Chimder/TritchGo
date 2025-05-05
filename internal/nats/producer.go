package nats

import (
	"context"
	"log"
	"log/slog"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type NatsProducer struct {
	Conn   *nats.Conn
	Stream jetstream.JetStream
}

func NewNatsProducer(ctx context.Context) *NatsProducer {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		slog.Info("Nats err conn")
		log.Fatal(err)
	}

	stream, err := jetstream.New(nc)
	if err != nil {
		slog.Info("Nats err jetstrean")
		log.Fatal(err)
	}

	_, err = stream.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:        "TRITCH_STATS",
		Description: "Msg for tritch stats",
		Subjects:    []string{"tritch.stats"},
		MaxBytes:    1024 * 1024 * 1024,
		Storage:     jetstream.FileStorage,
	})
	if err != nil {
		slog.Info("Nats err create stream")
		log.Fatal(err)
	}

	slog.Info("NATS JetStream producer running")
	return &NatsProducer{
		Conn:   nc,
		Stream: stream,
	}
}
func (np *NatsProducer) Publish(ctx context.Context, subject string, data []byte) error {
	_, err := np.Stream.Publish(ctx, subject, data)
	return err
}

func (np *NatsProducer) Close() {
	if np.Conn != nil && !np.Conn.IsClosed() {
		np.Conn.Close()
	}
}
