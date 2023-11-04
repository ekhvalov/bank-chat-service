package outbox

import (
	"context"
	"fmt"
	"time"

	"github.com/ekhvalov/bank-chat-service/internal/types"
)

func (s *Service) Put(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error) {
	jobID, err := s.repo.CreateJob(ctx, name, payload, availableAt)
	if err != nil {
		return types.JobIDNil, fmt.Errorf("create job: %v", err)
	}

	return jobID, nil
}
