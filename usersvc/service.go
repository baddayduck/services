package usersvc

import (
	"context"
	"errors"
	"sync"
)

var (
	// ErrAlreadyExists for when a user exists
	ErrAlreadyExists = errors.New("already exists")
	// ErrNotFound for when a user is not found
	ErrNotFound = errors.New("not found")
)

// Service is the user service itself
type Service interface {
	GetUser(ctx context.Context, id string) (User, error)
	AddUser(ctx context.Context, user User) error
	DeleteUser(ctx context.Context, id string) error
	HealthCheck() bool
}

// User represents a user
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type inmemService struct {
	mtx sync.RWMutex
	m   map[string]User
}

// NewInmemService creates an in memory service
func NewInmemService() Service {
	return &inmemService{
		m: map[string]User{},
	}
}

// GetUser get the user based on id
func (s *inmemService) GetUser(ctx context.Context, id string) (User, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	user, ok := s.m[id]
	if !ok {
		return User{}, ErrNotFound
	}
	return user, nil
}

// AddUser adds a user
func (s *inmemService) AddUser(ctx context.Context, user User) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if _, ok := s.m[user.ID]; ok {
		return ErrAlreadyExists
	}
	s.m[user.ID] = user
	return nil
}

// DeleteUser deletes a user
func (s *inmemService) DeleteUser(ctx context.Context, id string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if _, ok := s.m[id]; !ok {
		return ErrNotFound
	}
	delete(s.m, id)
	return nil
}

func (s *inmemService) HealthCheck() bool {
	// This really shouldn't always return true, but works for now.
	return true
}
