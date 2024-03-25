package kafka

import (
	"context"
	"fmt"

	"github.com/av-ugolkov/lingua-evo/internal/config"
	"github.com/segmentio/kafka-go"
)

type Writer struct {
	writer *kafka.Writer
}

func NewWriter(cfg *config.Config) *Writer {
	config := kafka.WriterConfig{
		Brokers: []string{cfg.Kafka.Addr()},
		Topic:   cfg.Kafka.Topic,
	}
	w := kafka.NewWriter(config)

	return &Writer{
		writer: w,
	}
}

func (w *Writer) SendMessage(ctx context.Context, key, value []byte) error {
	msg := kafka.Message{
		Key:   key,
		Value: value,
	}
	err := w.writer.WriteMessages(ctx, msg)
	if err != nil {
		return fmt.Errorf("delivery.kafka.Writer.SendMessage: %w", err)
	}

	return nil
}
