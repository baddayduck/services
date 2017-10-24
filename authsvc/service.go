package authsvc

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrNotFound      = errors.New("not found")
)

type Service interface {
	GetSession(ctx context.Context, id string) (Session, error)
	CreateSession(ctx context.Context, sess Session) error
	DeleteSession(ctx context.Context, id string) error
}

type Session struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type inmemService struct {
	mtx sync.RWMutex
	m   map[string]Session
}

func NewInmemService() Service {
	return &inmemService{
		m: map[string]Session{},
	}
}

func (s *inmemService) GetSession(ctx context.Context, id string) (Session, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	session, ok := s.m[id]
	if !ok {
		return Session{}, ErrNotFound
	}
	return session, nil
}

func (s *inmemService) CreateSession(ctx context.Context, sess Session) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if _, ok := s.m[sess.ID]; ok {
		return ErrAlreadyExists
	}
	s.m[sess.ID] = sess
	return nil
}

func (s *inmemService) DeleteSession(ctx context.Context, id string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if _, ok := s.m[id]; !ok {
		return ErrNotFound
	}
	delete(s.m, id)
	return nil
}
