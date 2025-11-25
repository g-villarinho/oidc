package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrSessionNotFound         = errors.New("session not found")
	ErrSessionExpired          = errors.New("session has expired")
	ErrInvalidSessionSignature = errors.New("invalid session signature")
)

type Session struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	ExpiresAt time.Time
	CreatedAt time.Time
}

func NewSession(userID uuid.UUID, ttl time.Duration) (*Session, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	return &Session{
		ID:        id,
		UserID:    userID,
		ExpiresAt: time.Now().Add(ttl),
		CreatedAt: time.Now(),
	}, nil
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

func (s *Session) TTL() time.Duration {
	return time.Until(s.ExpiresAt)
}
