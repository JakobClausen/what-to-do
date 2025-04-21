package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"what-to-do/internal/domain"
	"what-to-do/internal/store"
)

type Store struct {
	mu   sync.RWMutex
	next int64
	data map[int64]*domain.Command
}

var _ store.CommandStore = (*Store)(nil)

func New() *Store { return &Store{data: make(map[int64]*domain.Command)} }

func (s *Store) Create(ctx context.Context, c *domain.Command) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.next++
	c.ID = s.next
	now := time.Now()
	c.CreatedAt, c.UpdatedAt = now, now
	// store a copy to avoid accidental external mutation
	cloned := *c
	s.data[c.ID] = &cloned
	return c.ID, nil
}

func (s *Store) Get(ctx context.Context, id int64) (*domain.Command, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	c, ok := s.data[id]
	if !ok {
		return nil, errors.New("command not found")
	}
	clone := *c
	return &clone, nil
}

func (s *Store) List(ctx context.Context) ([]*domain.Command, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]*domain.Command, 0, len(s.data))
	for _, c := range s.data {
		copy := *c
		out = append(out, &copy)
	}
	return out, nil
}

func (s *Store) Update(ctx context.Context, c *domain.Command) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[c.ID]; !ok {
		return errors.New("command not found")
	}
	now := time.Now()
	c.UpdatedAt = now
	clone := *c
	s.data[c.ID] = &clone
	return nil
}

func (s *Store) Delete(ctx context.Context, id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[id]; !ok {
		return errors.New("command not found")
	}
	delete(s.data, id)
	return nil
}
