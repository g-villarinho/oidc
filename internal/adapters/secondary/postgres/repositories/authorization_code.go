package repositories

import (
	"context"

	"github.com/g-villarinho/oidc-server/internal/adapters/secondary/postgres/db"
	"github.com/g-villarinho/oidc-server/internal/core/domain"
	"github.com/g-villarinho/oidc-server/internal/core/ports"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthorizationCodeRepository struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

func NewAuthorizationCodeRepository(pool *pgxpool.Pool) ports.AuthorizationCodeRepository {
	return &AuthorizationCodeRepository{
		queries: db.New(pool),
		pool:    pool,
	}
}

func (r *AuthorizationCodeRepository) Create(ctx context.Context, code *domain.AuthorizationCode) error {
	userID := pgtype.UUID{
		Bytes: code.UserID,
		Valid: true,
	}
	expiresAt := pgtype.Timestamp{
		Time:  code.ExpiresAt,
		Valid: true,
	}

	codeChallenge := pgtype.Text{
		String: code.CodeChallenge,
		Valid:  true,
	}

	codeChallengeMethod := pgtype.Text{
		String: code.CodeChallengeMethod,
		Valid:  true,
	}

	_, err := r.queries.CreateAuthorizationCode(ctx, db.CreateAuthorizationCodeParams{
		Code:                code.Code,
		ClientID:            code.ClientID,
		UserID:              userID,
		RedirectUri:         code.RedirectURI,
		Scopes:              code.Scopes,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
		ExpiresAt:           expiresAt,
	})

	return err
}

func (a *AuthorizationCodeRepository) GetByCode(ctx context.Context, code string) (*domain.AuthorizationCode, error) {
	panic("unimplemented")
}

func (a *AuthorizationCodeRepository) Delete(ctx context.Context, code string) error {
	panic("unimplemented")
}
