package inmemmanagerpool

import (
	"context"
	"errors"
	"sync"

	managerpool "github.com/ekhvalov/bank-chat-service/internal/services/manager-pool"
	"github.com/ekhvalov/bank-chat-service/internal/types"
)

const (
	managersMax = 1000
)

func New() *Service {
	return &Service{mutex: &sync.Mutex{}, store: make(map[types.UserID]struct{})}
}

type Service struct {
	first *entry
	last  *entry
	total int
	mutex *sync.Mutex
	store map[types.UserID]struct{}
}

type entry struct {
	id   types.UserID
	next *entry
}

func (s *Service) Close() error {
	return nil
}

func (s *Service) Get(_ context.Context) (types.UserID, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if nil == s.first {
		return types.UserIDNil, managerpool.ErrNoAvailableManagers
	}

	e := s.first
	s.first = e.next
	delete(s.store, e.id)
	s.total--

	if e == s.last {
		s.last = nil
	}
	return e.id, nil
}

func (s *Service) Put(_ context.Context, managerID types.UserID) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, exists := s.store[managerID]
	if exists {
		return nil
	}

	if s.total == managersMax {
		return errors.New("managers limit reached")
	}

	e := entry{id: managerID}
	s.store[e.id] = struct{}{}
	s.total++

	if s.last != nil {
		s.last.next = &e
		s.last = &e
		return nil
	}

	s.last = &e
	s.first = &e
	return nil
}

func (s *Service) Contains(_ context.Context, managerID types.UserID) (bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, exists := s.store[managerID]
	return exists, nil
}

func (s *Service) Size() int {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.total
}
