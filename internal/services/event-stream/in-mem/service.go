package inmemeventstream

import (
	"context"
	"fmt"
	"sync"

	eventstream "github.com/ekhvalov/bank-chat-service/internal/services/event-stream"
	"github.com/ekhvalov/bank-chat-service/internal/types"
)

type Service struct {
	subs  map[types.UserID]*userSubscribers
	mutex *sync.Mutex
}

func New() *Service {
	return &Service{
		subs:  make(map[types.UserID]*userSubscribers),
		mutex: &sync.Mutex{},
	}
}

func (s *Service) Subscribe(ctx context.Context, userID types.UserID) (<-chan eventstream.Event, error) {
	us, ok := s.userSubscribers(userID)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if !ok {
		us = newUserSubscribers()
		s.subs[userID] = us
	}
	return us.subscribe(ctx), nil
}

func (s *Service) Publish(_ context.Context, userID types.UserID, event eventstream.Event) error {
	if err := event.Validate(); err != nil {
		return fmt.Errorf("invalid event: %v", err)
	}

	us, ok := s.userSubscribers(userID)
	if !ok {
		return nil
	}

	us.publish(event)

	return nil
}

func (s *Service) userSubscribers(userID types.UserID) (*userSubscribers, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	us, ok := s.subs[userID]
	return us, ok
}

func (s *Service) Close() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	userIDs := make([]types.UserID, 0, len(s.subs))
	for userID, subscribers := range s.subs {
		subscribers.close()
		userIDs = append(userIDs, userID)
	}
	for _, uid := range userIDs {
		delete(s.subs, uid)
	}
	return nil
}

func newUserSubscribers() *userSubscribers {
	return &userSubscribers{
		subscribers: make([]*subscriber, 0),
		mutex:       &sync.Mutex{},
		wg:          &sync.WaitGroup{},
	}
}

type userSubscribers struct {
	subscribers []*subscriber
	cancels     []context.CancelFunc
	mutex       *sync.Mutex
	wg          *sync.WaitGroup
}

func (s *userSubscribers) close() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, cancel := range s.cancels {
		cancel()
	}
	s.wg.Wait()
}

func (s *userSubscribers) subscribe(ctx context.Context) <-chan eventstream.Event {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	subCtx, cancel := context.WithCancel(ctx)
	sub := newSubscriber(subCtx)
	s.wg.Add(1)
	go func() {
		sub.subscribe()
		s.wg.Done()
	}()
	s.subscribers = append(s.subscribers, sub)
	s.cancels = append(s.cancels, cancel)
	return sub.subscription
}

func (s *userSubscribers) publish(event eventstream.Event) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, sub := range s.subscribers {
		sub.publish(event)
	}
}

type event struct {
	value eventstream.Event
	next  *event
}

func newSubscriber(ctx context.Context) *subscriber {
	return &subscriber{
		ctx:          ctx,
		mutex:        &sync.Mutex{},
		hasEvent:     make(chan struct{}, 1), // buffered
		subscription: make(chan eventstream.Event),
	}
}

type subscriber struct {
	ctx          context.Context
	first        *event
	last         *event
	mutex        *sync.Mutex
	hasEvent     chan struct{}
	subscription chan eventstream.Event
}

func (s *subscriber) subscribe() {
	defer func() {
		close(s.hasEvent)
		close(s.subscription)
	}()
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-s.hasEvent:
			s.mutex.Lock()
			var e *event
			if s.first != nil {
				e = s.first
				if s.first == s.last {
					s.first = nil
					s.last = nil
				} else {
					s.first = s.first.next
				}
				if s.first != nil {
					select {
					case s.hasEvent <- struct{}{}:
					default:
					}
				}
			}
			s.mutex.Unlock()
			if e != nil {
				select {
				case s.subscription <- e.value:
				case <-s.ctx.Done():
					return
				}
			}
		}
	}
}

func (s *subscriber) publish(value eventstream.Event) {
	select {
	case <-s.ctx.Done():
		return
	default:
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	e := &event{value: value}
	if nil == s.last {
		s.first = e
		s.last = e
	} else {
		s.last.next = e
		s.last = e
	}
	select {
	case s.hasEvent <- struct{}{}:
	default:
	}
}
