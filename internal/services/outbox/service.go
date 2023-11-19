package outbox

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	jobsrepo "github.com/ekhvalov/bank-chat-service/internal/repositories/jobs"
	"github.com/ekhvalov/bank-chat-service/internal/types"
)

type jobsRepository interface {
	CreateJob(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error)
	FindAndReserveJob(ctx context.Context, until time.Time) (jobsrepo.Job, error)
	CreateFailedJob(ctx context.Context, name, payload, reason string) error
	DeleteJob(ctx context.Context, jobID types.JobID) error
}

type transactor interface {
	RunInTx(ctx context.Context, f func(context.Context) error) error
}

//go:generate options-gen -out-filename=service_options.gen.go -from-struct=Options
type Options struct {
	workers    int            `option:"mandatory" validate:"min=1,max=32"`
	idleTime   time.Duration  `option:"mandatory" validate:"min=100ms,max=10s"`
	reserveFor time.Duration  `option:"mandatory" validate:"min=1s,max=10m"`
	repo       jobsRepository `option:"mandatory" validate:"required"`
	db         transactor     `option:"mandatory" validate:"required"`
	lg         *zap.Logger    `option:"mandatory" validate:"required"`
}

type Service struct {
	Options
	registry map[string]Job
	mutex    *sync.Mutex
}

func New(opts Options) (*Service, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate options: %v", err)
	}
	return &Service{
		Options:  opts,
		registry: make(map[string]Job),
		mutex:    &sync.Mutex{},
	}, nil
}

func (s *Service) RegisterJob(job Job) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, exists := s.registry[job.Name()]
	if exists {
		return fmt.Errorf("job %q is already registered", job.Name())
	}

	s.registry[job.Name()] = job

	return nil
}

func (s *Service) MustRegisterJob(job Job) {
	if err := s.RegisterJob(job); err != nil {
		panic(fmt.Errorf("register job: %v", err))
	}
}

func (s *Service) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)
	for i := 0; i < s.workers; i++ {
		eg.Go(func() error { return s.handleJobs(ctx) })
	}
	if err := eg.Wait(); err != nil && errors.Is(err, context.Canceled) {
		return fmt.Errorf("run: %v", err)
	}
	return nil
}

func (s *Service) handleJobs(ctx context.Context) error {
	for nil == ctx.Err() {
		job, err := s.repo.FindAndReserveJob(ctx, time.Now().Add(s.reserveFor))
		if err != nil {
			if errors.Is(err, jobsrepo.ErrNoJobs) {
				idle(ctx, s.idleTime)
				continue
			}
			return fmt.Errorf("find and reserve: %v", err)
		}
		if err = s.handleJob(ctx, job); err != nil {
			return fmt.Errorf("handle job %q: %v", job.ID, err)
		}
	}
	return nil
}

func (s *Service) handleJob(ctx context.Context, job jobsrepo.Job) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	handler, ok := s.registry[job.Name]
	if !ok {
		if err := s.moveJobToFailed(ctx, job, "handler is not found"); err != nil {
			return fmt.Errorf("fail no handler job: %v", err)
		}
		return nil
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, handler.ExecutionTimeout())
	defer cancel()
	err := handler.Handle(ctxTimeout, job.Payload)
	if nil == err {
		if err := s.repo.DeleteJob(ctx, job.ID); err != nil {
			return fmt.Errorf("delete job: %v", err)
		}
		return nil
	}

	s.lg.Error("handle payload", zap.String("job_id", job.ID.String()), zap.Error(err))
	if job.Attempts >= handler.MaxAttempts() {
		if err := s.moveJobToFailed(ctx, job, "attempts limit exceeded"); err != nil {
			return fmt.Errorf("fail attempts limit exceeded job: %v", err)
		}
	}
	return nil
}

func (s *Service) moveJobToFailed(ctx context.Context, job jobsrepo.Job, reason string) error {
	return s.db.RunInTx(ctx, func(ctx context.Context) error {
		if err := s.repo.DeleteJob(ctx, job.ID); err != nil {
			return fmt.Errorf("delete job %q: %v", job.ID, err)
		}
		if err := s.repo.CreateFailedJob(ctx, job.Name, job.Payload, reason); err != nil {
			return fmt.Errorf("create failed job: %v", err)
		}
		return nil
	})
}

func idle(ctx context.Context, idleTime time.Duration) {
	timer := time.NewTimer(idleTime)
	defer timer.Stop()

	select {
	case <-timer.C:
	case <-ctx.Done():
	}
}
