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

	nonce := pgtype.Text{
		String: code.Nonce,
		Valid:  code.Nonce != "",
	}

	codeChallenge := pgtype.Text{
		String: code.CodeChallenge,
		Valid:  true,
	}

	codeChallengeMethod := pgtype.Text{
		String: code.CodeChallengeMethod,
		Valid:  true,
	}

	scopes := code.Scopes
	if scopes == nil {
		scopes = []string{}
	}

	_, err := r.queries.CreateAuthorizationCode(ctx, db.CreateAuthorizationCodeParams{
		Code:                code.Code,
		ClientID:            code.ClientID,
		UserID:              userID,
		RedirectUri:         code.RedirectURI,
		Scopes:              scopes,
		Nonce:               nonce,
		CodeChallenge:       codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
		ExpiresAt:           expiresAt,
	})

	return err
}

func (a *AuthorizationCodeRepository) GetByCode(ctx context.Context, code string) (*domain.AuthorizationCode, error) {
	ac, err := a.queries.GetAuthorizationCode(ctx, code)
	if err != nil {
		if isNotFound(err) {
			return nil, ports.ErrNotFound
		}

		return nil, err
	}

	return &domain.AuthorizationCode{
		Code:                ac.Code,
		ClientID:            ac.ClientID,
		UserID:              ac.UserID.Bytes,
		RedirectURI:         ac.RedirectUri,
		Scopes:              ac.Scopes,
		Nonce:               ac.Nonce.String,
		CodeChallenge:       ac.CodeChallenge.String,
		CodeChallengeMethod: ac.CodeChallengeMethod.String,
		Used:                ac.Used,
		ExpiresAt:           ac.ExpiresAt.Time,
		CreatedAt:           ac.CreatedAt.Time,
	}, nil
}

func (a *AuthorizationCodeRepository) Delete(ctx context.Context, code string) error {
	return a.queries.DeleteAuthorizationCode(ctx, code)
}

func (r *AuthorizationCodeRepository) MarkAsUsed(ctx context.Context, code string) error {
	return r.queries.MarkAuthorizationCodeAsUsed(ctx, code)
}
