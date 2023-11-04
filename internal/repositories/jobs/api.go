package jobsrepo

import (
	"context"
	"fmt"
	"time"

	store "github.com/ekhvalov/bank-chat-service/internal/store/gen"
	"github.com/ekhvalov/bank-chat-service/internal/store/gen/job"
	"github.com/ekhvalov/bank-chat-service/internal/types"
)

type Job struct {
	ID       types.JobID
	Name     string
	Payload  string
	Attempts int
}

func (r *Repo) FindAndReserveJob(ctx context.Context, until time.Time) (Job, error) {
	var j *store.Job

	findAndReserve := func(ctx context.Context) error {
		tx := store.TxFromContext(ctx)
		if nil == tx {
			return fmt.Errorf("no transaction in context")
		}

		now := time.Now()
		var err error
		j, err = tx.Job.Query().
			Where(
				job.AvailableAtLTE(now),
				job.ReservedUntilLTE(now),
			).
			ForUpdate().
			First(ctx)
		if err != nil {
			if store.IsNotFound(err) {
				return ErrNoJobs
			}
			return fmt.Errorf("find job: %v", err)
		}

		j, err = j.Update().AddAttempts(1).SetReservedUntil(until).Save(ctx)
		if err != nil {
			if isAttemptsExceededError(err) {
				return ErrAttemptsExceeded
			}
			return fmt.Errorf("update job: %v", err)
		}
		return nil
	}
	if err := r.db.RunInTx(ctx, findAndReserve); err != nil {
		return Job{}, err
	}

	return Job{
		ID:       j.ID,
		Name:     j.Name,
		Payload:  j.Payload,
		Attempts: j.Attempts,
	}, nil
}

func (r *Repo) CreateJob(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error) {
	j, err := r.db.Job(ctx).
		Create().
		SetName(name).
		SetPayload(payload).
		SetAvailableAt(availableAt).
		Save(ctx)
	if err != nil {
		return types.JobIDNil, fmt.Errorf("create job: %v", err)
	}
	return j.ID, nil
}

func (r *Repo) CreateFailedJob(ctx context.Context, name, payload, reason string) error {
	_, err := r.db.FailedJob(ctx).
		Create().
		SetName(name).
		SetPayload(payload).
		SetReason(reason).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("create failedJob: %v", err)
	}
	return nil
}

func (r *Repo) DeleteJob(ctx context.Context, jobID types.JobID) error {
	count, err := r.db.Job(ctx).Delete().Where(job.ID(jobID)).Exec(ctx)
	if err != nil {
		return fmt.Errorf("delete job %q: %v", jobID, err)
	}
	if count == 0 {
		return ErrNoJobs
	}
	return nil
}
