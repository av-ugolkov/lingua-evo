package analytics

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/av-ugolkov/lingua-evo/internal/delivery/kafka"

	"github.com/google/uuid"
)

var (
	topic = []byte("user_action")
)

type Analytic struct {
	writer *kafka.Writer
}

var instant Analytic

func SetKafka(w *kafka.Writer) {
	instant.writer = w
}

func SendToKafka(uid uuid.UUID, action string) {
	if instant.writer == nil {
		return
	}
	go func() {
		b, err := json.Marshal(Action{
			UserID:    uid,
			Action:    action,
			CreatedAt: time.Now().UTC(),
		})
		if err != nil {
			slog.Error(fmt.Errorf("pkg.analytics.SendToKafka - marshal: %v", err).Error())
			return
		}
		err = instant.writer.SendMessage(context.Background(), topic, b)
		if err != nil {
			slog.Error(fmt.Errorf("pkg.analytics.SendToKafka: %v", err).Error())
		}
	}()
}
