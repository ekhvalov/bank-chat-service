package inmemmanagerpool

import (
	"context"
	"errors"
	"sync"

	managerpool "github.com/ekhvalov/bank-chat-service/internal/services/manager-pool"
	"github.com/ekhvalov/bank-chat-service/internal/types"
)

const (
	serviceName = "manager-pool"
	managersMax = 1000
)

func New() *Service {
	return &Service{mutex: &sync.RWMutex{}, store: make(map[types.UserID]struct{})}
}

type Service struct {
	first *entry
	last  *entry
	total int
	mutex *sync.RWMutex
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
	s.mutex.RLock()
	if nil == s.first {
		s.mutex.RUnlock()
		return types.UserIDNil, managerpool.ErrNoAvailableManagers
	}
	s.mutex.RUnlock()
	s.mutex.Lock()
	defer s.mutex.Unlock()
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
	s.mutex.RLock()
	_, exists := s.store[managerID]
	s.mutex.RUnlock()
	if exists {
		return nil
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
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
	s.mutex.RLock()
	_, exists := s.store[managerID]
	s.mutex.RUnlock()
	return exists, nil
}

func (s *Service) Size() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.total
}
