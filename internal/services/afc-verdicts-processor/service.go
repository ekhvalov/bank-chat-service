package afcverdictsprocessor

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/ekhvalov/bank-chat-service/internal/types"
	"github.com/ekhvalov/bank-chat-service/internal/validator"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/service_mock.gen.go -typed -package=afcverdictsprocessormocks

const serviceName = "afc-verdicts-processor"

type messagesRepository interface {
	MarkAsVisibleForManager(ctx context.Context, msgID types.MessageID) error
	BlockMessage(ctx context.Context, msgID types.MessageID) error
}

type outboxService interface {
	Put(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error)
}

type transactor interface {
	RunInTx(ctx context.Context, f func(context.Context) error) error
}

//go:generate options-gen -out-filename=service_options.gen.go -from-struct=Options
type Options struct {
	processBatchSize int `default:"1" validate:"min=1,max=32"`

	backoffInitialInterval time.Duration `default:"100ms" validate:"min=50ms,max=1s"`
	backoffMaxElapsedTime  time.Duration `default:"5s" validate:"min=500ms,max=1m"`

	brokers         []string `option:"mandatory" validate:"min=1"`
	consumers       int      `option:"mandatory" validate:"min=1,max=16"`
	consumerGroup   string   `option:"mandatory" validate:"required"`
	verdictsTopic   string   `option:"mandatory" validate:"required"`
	verdictsSignKey string

	readerFactory KafkaReaderFactory `option:"mandatory" validate:"required"`
	dlqWriter     KafkaDLQWriter     `option:"mandatory" validate:"required"`

	txtor   transactor         `option:"mandatory" validate:"required"`
	msgRepo messagesRepository `option:"mandatory" validate:"required"`
	outBox  outboxService      `option:"mandatory" validate:"required"`
	logger  *zap.Logger        `option:"mandatory" validate:"required"`
}

type Service struct {
	Options
	kafkaReader KafkaReader
	rsaPubKey   *rsa.PublicKey
	processor   *processor
}

func New(opts Options) (*Service, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate options: %v", err)
	}
	var pk *rsa.PublicKey
	var err error

	if opts.verdictsSignKey != "" {
		pk, err = jwt.ParseRSAPublicKeyFromPEM([]byte(opts.verdictsSignKey))
		if err != nil {
			return nil, fmt.Errorf("parse rsa pubkey: %v", err)
		}
	}

	return &Service{
		Options:     opts,
		kafkaReader: opts.readerFactory(opts.brokers, opts.consumerGroup, opts.verdictsTopic),
		rsaPubKey:   pk,
		processor: &processor{
			txtor:                  opts.txtor,
			msgRepo:                opts.msgRepo,
			outBox:                 opts.outBox,
			backoffInitialInterval: opts.backoffInitialInterval,
			backoffMaxElapsedTime:  opts.backoffMaxElapsedTime,
		},
	}, nil
}

func (s *Service) Run(ctx context.Context) (err error) {
	defer func() {
		if errClose := s.dlqWriter.Close(); errClose != nil {
			err = errors.Join(err, fmt.Errorf("DLQ writer close: %v", errClose))
		}
		if errClose := s.kafkaReader.Close(); errClose != nil {
			err = errors.Join(err, fmt.Errorf("kafka reader close: %v", errClose))
		}
		s.logger.Debug("done")
	}()

	dlqCh := make(chan kafka.Message)
	commitCh := make(chan kafka.Message)
	resultsChs := make([]<-chan result, 0, s.consumers)

	eg, ctx := errgroup.WithContext(ctx)

	for i := 0; i < s.consumers; i++ {
		logger := s.logger.With(zap.Int("worker", i+1))
		resCh, errCh := s.consume(ctx, logger)
		resultsChs = append(resultsChs, resCh)
		eg.Go(func() error {
			return <-errCh
		})
	}

	wg := &sync.WaitGroup{}
	for _, resCh := range resultsChs {
		wg.Add(1)
		go func(resCh <-chan result) {
			for res := range resCh {
				if res.err != nil {
					dlqCh <- dlqMessage(res.msg, res.err)
				}
				commitCh <- res.msg
			}
			wg.Done()
		}(resCh)
	}
	go func() {
		wg.Wait()
		close(dlqCh)
		close(commitCh)
	}()

	eg.Go(func() error {
		if err := s.dlqLoop(dlqCh); err != nil {
			return fmt.Errorf("DLQ loop: %v", err)
		}
		return nil
	})

	eg.Go(func() error {
		if err := s.commitLoop(commitCh); err != nil {
			return fmt.Errorf("commit loop: %v", err)
		}
		return nil
	})

	err = eg.Wait()
	if err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}

type result struct {
	msg kafka.Message
	err error
}

func (s *Service) consume(ctx context.Context, log *zap.Logger) (<-chan result, <-chan error) {
	resCh := make(chan result)
	errCh := make(chan error, 1)

	go func() {
		defer close(resCh)
		defer close(errCh)

		for {
			msg, err := s.kafkaReader.FetchMessage(ctx)
			if err != nil {
				if !(errors.Is(err, io.EOF) || errors.Is(err, context.Canceled)) {
					log.Error("fetch message error", zap.Error(err))
					errCh <- fmt.Errorf("fetch message: %v", err)
				}
				return
			}

			var v verdict
			if v, err = s.parseVerdict(msg.Value); err != nil {
				log.Warn("parse verdict", zap.Error(err))
				select {
				case resCh <- result{msg: msg, err: err}:
				case <-ctx.Done():
					return
				}
				continue
			}

			err = s.processor.process(ctx, v)
			if err != nil {
				log.Warn(
					"process error",
					zap.Error(err),
					zap.Stringer("message_id", v.MessageID),
					zap.Stringer("chat_id", v.ChatID),
				)
			}
			select {
			case resCh <- result{msg: msg, err: err}:
			case <-ctx.Done():
				return
			}
		}
	}()

	return resCh, errCh
}

func (s *Service) parseVerdict(data []byte) (verdict, error) {
	var v verdict
	if s.rsaPubKey != nil {
		parser := jwt.Parser{ValidMethods: []string{jwt.SigningMethodRS256.Alg()}}
		t, err := parser.Parse(string(data), func(token *jwt.Token) (interface{}, error) {
			return s.rsaPubKey, nil
		})
		if err != nil {
			return verdict{}, fmt.Errorf("jwt parse: %v", err)
		}
		parts := strings.Split(t.Raw, ".")
		data, _ = jwt.DecodeSegment(parts[1])
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return v, fmt.Errorf("unmarshal: %v", err)
	}

	if err := validator.Validator.Struct(&v); err != nil {
		return v, fmt.Errorf("validate: %v", err)
	}

	return v, nil
}

func (s *Service) dlqLoop(msgCh <-chan kafka.Message) error {
	messages := make([]kafka.Message, s.processBatchSize)
	for {
		count := fillSliceFromChan(messages, msgCh)
		if count == 0 {
			return nil
		}

		if err := s.dlqWriter.WriteMessages(context.Background(), messages[:count]...); err != nil {
			s.logger.Warn(
				"write to DLQ error",
				zap.Error(err),
				zap.Dict("messages", messagesToLogFields(messages[:count])...),
			)
			return fmt.Errorf("write to DLQ: %v", err)
		}
	}
}

func (s *Service) commitLoop(msgCh <-chan kafka.Message) error {
	messages := make([]kafka.Message, s.consumers)
	for {
		count := fillSliceFromChan(messages, msgCh)
		if count == 0 {
			return nil
		}

		if err := s.kafkaReader.CommitMessages(context.Background(), messages[:count]...); err != nil {
			return fmt.Errorf("commig messages: %v", err)
		}
	}
}

func fillSliceFromChan(messages []kafka.Message, msgCh <-chan kafka.Message) int {
	count := 0
	var ok bool

	messages[count], ok = <-msgCh
	if !ok {
		return 0
	}
	count++

	for count < len(messages) {
		select {
		case messages[count], ok = <-msgCh:
			if !ok {
				return count
			}
			count++

		default:
			return count
		}
	}

	return count
}

func messagesToLogFields(messages []kafka.Message) []zap.Field {
	fields := make([]zap.Field, len(messages))
	for i, message := range messages {
		headers := make([]zap.Field, len(message.Headers))
		for j, header := range message.Headers {
			headers[j] = zap.Binary(header.Key, header.Value)
		}
		fields[i] = zap.Dict("message",
			zap.Int("partition", message.Partition),
			zap.Dict("headers", headers...),
			zap.Int64("offset", message.Offset),
		)
	}
	return fields
}

func dlqMessage(msg kafka.Message, err error) kafka.Message {
	dlqMsg := kafka.Message{Value: msg.Value}
	dlqMsg.Headers = append(dlqMsg.Headers, kafka.Header{Key: "LAST_ERROR", Value: []byte(err.Error())})
	dlqMsg.Headers = append(dlqMsg.Headers, kafka.Header{Key: "ORIGINAL_PARTITION", Value: []byte(strconv.Itoa(msg.Partition))})
	return dlqMsg
}
