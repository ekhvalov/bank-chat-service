package outbox

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/ekhvalov/bank-chat-service/internal/types"
)

func (s *Service) Put(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error) {
	jobID, err := s.jobsRepo.CreateJob(ctx, name, payload, availableAt)
	if err != nil {
		return types.JobIDNil, fmt.Errorf("create job: %v", err)
	}
	s.lg.Debug("put job", zap.String("name", name), zap.String("payload", payload))
	return jobID, nil
}
