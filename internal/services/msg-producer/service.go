package msgproducer

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type KafkaWriter interface {
	io.Closer
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
}

//go:generate options-gen -out-filename=service_options.gen.go -from-struct=Options -defaults-from=func
type Options struct {
	wr           KafkaWriter `option:"mandatory" validate:"required"`
	logger       *zap.Logger `option:"mandatory" validate:"required"`
	encryptKey   string      `validate:"omitempty,hexadecimal"`
	nonceFactory func(size int) ([]byte, error)
}

func getDefaultOptions() Options {
	return Options{
		nonceFactory: defaultNonceFactory,
	}
}

type Service struct {
	wr        KafkaWriter
	encryptor encryptor
	logger    *zap.Logger
}

func New(opts Options) (*Service, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate options: %v", err)
	}

	var err error
	var enc encryptor = &encryptorPlain{}
	if len(opts.encryptKey) > 0 {
		enc, err = newEncryptorAEAD(opts.encryptKey, opts.nonceFactory)
		if err != nil {
			return nil, fmt.Errorf("create encryptor: %v", err)
		}
	} else {
		opts.logger.Info("encryption disabled")
	}

	return &Service{
		wr:        opts.wr,
		encryptor: enc,
		logger:    opts.logger,
	}, nil
}

func defaultNonceFactory(size int) ([]byte, error) {
	nonce := make([]byte, size)
	_, err := rand.Read(nonce)
	return nonce, err
}
