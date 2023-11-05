package msgproducer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"

	"github.com/ekhvalov/bank-chat-service/internal/types"
)

type Message struct {
	ID         types.MessageID `json:"id"`
	ChatID     types.ChatID    `json:"chatId"`
	Body       string          `json:"body"`
	FromClient bool            `json:"fromClient"`
}

func (s *Service) ProduceMessage(ctx context.Context, msg Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("json marshal: %v", err)
	}
	encData, err := s.encryptor.encrypt(data)
	if err != nil {
		return fmt.Errorf("encrypt: %v", err)
	}
	err = s.wr.WriteMessages(ctx, kafka.Message{
		Key:   []byte(msg.ChatID.String()),
		Value: encData,
	})
	if err != nil {
		return fmt.Errorf("write to kafka: %v", err)
	}
	return nil
}

func (s *Service) Close() error {
	if err := s.wr.Close(); err != nil {
		return fmt.Errorf("close kafka writer: %v", err)
	}
	return nil
}
