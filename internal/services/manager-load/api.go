package managerload

import (
	"context"
	"fmt"

	"github.com/ekhvalov/bank-chat-service/internal/types"
)

func (s *Service) CanManagerTakeProblem(ctx context.Context, managerID types.UserID) (bool, error) {
	count, err := s.problemsRepo.CountManagerOpenProblems(ctx, managerID)
	if err != nil {
		return false, fmt.Errorf("count manager open problems: %v", err)
	}
	return count < s.maxProblemsAtTime, nil
}
