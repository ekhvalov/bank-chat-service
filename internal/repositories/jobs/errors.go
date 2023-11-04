package jobsrepo

import (
	"errors"

	"github.com/ekhvalov/bank-chat-service/internal/store/gen"
)

var (
	ErrNoJobs           = errors.New("no jobs found")
	ErrAttemptsExceeded = errors.New("attempts limit exceeded")
)

func isAttemptsExceededError(err error) bool {
	if errVal := new(gen.ValidationError); errors.As(err, &errVal) {
		return errVal.Name == "attempts"
	}
	return false
}
