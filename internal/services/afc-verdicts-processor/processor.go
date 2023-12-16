package afcverdictsprocessor

import (
	"context"
	"fmt"
	"time"

	clientmessageblockedjob "github.com/ekhvalov/bank-chat-service/internal/services/outbox/jobs/client-message-blocked"
	clientmessagesentjob "github.com/ekhvalov/bank-chat-service/internal/services/outbox/jobs/client-message-sent"
	"github.com/ekhvalov/bank-chat-service/internal/types"
)

type processor struct {
	txtor   transactor         `option:"mandatory" validate:"required"`
	msgRepo messagesRepository `option:"mandatory" validate:"required"`
	outBox  outboxService      `option:"mandatory" validate:"required"`

	backoffInitialInterval time.Duration `default:"100ms" validate:"min=50ms,max=1s"`
	backoffMaxElapsedTime  time.Duration `default:"5s" validate:"min=500ms,max=1m"`
}

func (p *processor) process(ctx context.Context, v verdict) error {
	var handler func(ctx context.Context, id types.MessageID) error

	switch v.Status {
	case StatusOK:
		handler = p.handleStatusOK
	case StatusSuspicious:
		handler = p.handleStatusSuspicious
	default:
		return fmt.Errorf("unknown verdict status: %s", v.Status)
	}

	return p.runWithBackoff(ctx, func() error {
		return handler(ctx, v.MessageID)
	})
}

func (p *processor) handleStatusOK(ctx context.Context, id types.MessageID) error {
	return p.txtor.RunInTx(ctx, func(ctx context.Context) error {
		err := p.msgRepo.MarkAsVisibleForManager(ctx, id)
		if err != nil {
			return fmt.Errorf("mark as visible for manager: %v", err)
		}

		_, err = p.outBox.Put(ctx, clientmessagesentjob.Name, id.String(), time.Now())
		if err != nil {
			return fmt.Errorf("outbox put %s: %v", clientmessagesentjob.Name, err)
		}

		return nil
	})
}

func (p *processor) handleStatusSuspicious(ctx context.Context, id types.MessageID) error {
	return p.txtor.RunInTx(ctx, func(ctx context.Context) error {
		err := p.msgRepo.BlockMessage(ctx, id)
		if err != nil {
			return fmt.Errorf("block message: %v", err)
		}

		_, err = p.outBox.Put(ctx, clientmessageblockedjob.Name, id.String(), time.Now())
		if err != nil {
			return fmt.Errorf("outbox put %s: %v", clientmessagesentjob.Name, err)
		}

		return nil
	})
}

func (p *processor) runWithBackoff(ctx context.Context, workload func() error) error {
	sleepDuration := p.backoffInitialInterval
	timeout := time.NewTimer(p.backoffMaxElapsedTime)
	defer timeout.Stop()

	var err error
	for {
		if err = workload(); nil == err {
			return nil
		}
		sleep := time.NewTimer(sleepDuration)
		select {
		case <-timeout.C:
			sleep.Stop()
			return err
		case <-ctx.Done():
			sleep.Stop()
			return err
		case <-sleep.C:
			sleepDuration *= 2
		}
	}
}
