package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type Writer struct {
	writer *kafka.Writer
}

func NewWriter(addr, topic string) *Writer {
	config := kafka.WriterConfig{
		Brokers: []string{addr},
		Topic:   topic,
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
