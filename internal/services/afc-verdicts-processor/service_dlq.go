package afcverdictsprocessor

import (
	"context"
	"io"

	"github.com/segmentio/kafka-go"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/dlq_writer_mock.gen.go -typed -package=afcverdictsprocessormocks

type KafkaDLQWriter interface {
	io.Closer
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
}

func NewKafkaDLQWriter(brokers []string, topic string) KafkaDLQWriter {
	return &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Async:        false,
		RequiredAcks: kafka.RequireOne,
	}
}