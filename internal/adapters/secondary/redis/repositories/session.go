package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/g-villarinho/oidc-server/internal/core/domain"
	"github.com/g-villarinho/oidc-server/internal/core/ports"
	"github.com/g-villarinho/oidc-server/pkg/cache"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type SessionRepository struct {
	client *redis.Client
}

func NewSessionRepository(client *redis.Client) ports.SessionRepository {
	return &SessionRepository{
		client: client,
	}
}

func (r *SessionRepository) Create(ctx context.Context, session *domain.Session) error {
	key := cache.SessionKey(session.ID.String())

	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("marshal session: %w", err)
	}

	ttl := time.Until(session.ExpiresAt)
	if err := r.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("store session: %w", err)
	}

	return nil

}

func (r *SessionRepository) GetByID(ctx context.Context, sessionID uuid.UUID) (*domain.Session, error) {
	key := cache.SessionKey(sessionID.String())

	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, ports.ErrNotFound
		}
		return nil, fmt.Errorf("get session: %w", err)
	}

	var session domain.Session
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return nil, fmt.Errorf("unmarshal session: %w", err)
	}

	return &session, nil
}

func (r *SessionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Session, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *SessionRepository) Delete(ctx context.Context, sessionID uuid.UUID) error {
	key := cache.SessionKey(sessionID.String())

	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}

	return nil
}
