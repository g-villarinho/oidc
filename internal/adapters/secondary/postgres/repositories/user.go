package repositories

import (
	"context"
	"fmt"

	"github.com/g-villarinho/oidc-server/internal/adapters/secondary/postgres/db"
	"github.com/g-villarinho/oidc-server/internal/core/domain"
	"github.com/g-villarinho/oidc-server/internal/core/ports"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) ports.UserRepository {
	return &UserRepository{
		queries: db.New(pool),
		pool:    pool,
	}
}

func (u *UserRepository) Create(ctx context.Context, user *domain.User) error {
	pgUUID := pgtype.UUID{
		Bytes: user.ID,
		Valid: true,
	}

	_, err := u.queries.CreateUser(ctx, db.CreateUserParams{
		ID:            pgUUID,
		Name:          user.Name,
		Email:         user.Email,
		PasswordHash:  user.PasswordHash,
		EmailVerified: user.EmailVerified,
	})

	if err != nil {
		return fmt.Errorf("persist user: %w", err)
	}

	return nil
}

func (u *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := u.queries.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return u.toDomain(user), nil
}

func (u *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	pgUUID := pgtype.UUID{
		Bytes: id,
		Valid: true,
	}

	user, err := u.queries.GetByID(ctx, pgUUID)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return u.toDomain(user), nil
}

func (u *UserRepository) Update(ctx context.Context, user *domain.User) error {
	pgUUID := pgtype.UUID{
		Bytes: user.ID,
		Valid: true,
	}

	_, err := u.queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:            pgUUID,
		Name:          user.Name,
		Email:         user.Email,
		PasswordHash:  user.PasswordHash,
		EmailVerified: user.EmailVerified,
	})

	if err != nil {
		if isUniqueViolation(err) {
			return ports.ErrUniqueKeyViolation
		}

		return fmt.Errorf("update user: %w", err)
	}

	return nil
}

func (r *UserRepository) toDomain(user db.User) *domain.User {
	return &domain.User{
		ID:            user.ID.Bytes,
		Email:         user.Email,
		PasswordHash:  user.PasswordHash,
		Name:          user.Name,
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt.Time,
		UpdatedAt:     user.UpdatedAt.Time,
	}
}
