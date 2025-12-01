package repositories

import (
	"context"
	"time"

	"github.com/g-villarinho/oidc-server/internal/adapters/secondary/postgres/db"
	"github.com/g-villarinho/oidc-server/internal/core/domain"
	"github.com/g-villarinho/oidc-server/internal/core/ports"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenRepository struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

func NewTokenRepository(pool *pgxpool.Pool) ports.TokenRepository {
	return &TokenRepository{
		queries: db.New(pool),
		pool:    pool,
	}
}

func (r *TokenRepository) Create(ctx context.Context, token *domain.Token) error {
	id := pgtype.UUID{
		Bytes: token.ID,
		Valid: true,
	}

	userID := pgtype.UUID{
		Bytes: token.UserID,
		Valid: true,
	}

	accessTokenExpiresAt := pgtype.Timestamp{
		Time:  token.AccessTokenExpiresAt,
		Valid: true,
	}

	refreshTokenExpiresAt := pgtype.Timestamp{
		Time:  token.RefreshTokenExpiresAt,
		Valid: true,
	}

	var authCodePtr pgtype.Text
	if token.AuthorizationCode != nil {
		authCodePtr = pgtype.Text{
			String: *token.AuthorizationCode,
			Valid:  true,
		}
	}

	scopes := token.Scopes
	if scopes == nil {
		scopes = []string{}
	}

	_, err := r.queries.CreateToken(ctx, db.CreateTokenParams{
		ID:                   id,
		AccessTokenHash:      token.AccessTokenHash,
		RefreshTokenHash:     token.RefreshTokenHash,
		AuthorizationCode:    authCodePtr,
		ClientID:             token.ClientID,
		UserID:               userID,
		Scopes:               scopes,
		TokenType:            token.TokenType,
		AccessTokenExpiresAt: accessTokenExpiresAt,
		RefreshTokenExpiresAt: refreshTokenExpiresAt,
	})

	return err
}

func (r *TokenRepository) GetByAccessTokenHash(ctx context.Context, accessTokenHash string) (*domain.Token, error) {
	t, err := r.queries.GetTokenByAccessTokenHash(ctx, accessTokenHash)
	if err != nil {
		if isNotFound(err) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}

	return r.mapTokenToDomain(t), nil
}

func (r *TokenRepository) GetByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (*domain.Token, error) {
	t, err := r.queries.GetTokenByRefreshTokenHash(ctx, refreshTokenHash)
	if err != nil {
		if isNotFound(err) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}

	return r.mapTokenToDomain(t), nil
}

func (r *TokenRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Token, error) {
	tokenID := pgtype.UUID{
		Bytes: id,
		Valid: true,
	}

	t, err := r.queries.GetTokenByID(ctx, tokenID)
	if err != nil {
		if isNotFound(err) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}

	return r.mapTokenToDomain(t), nil
}

func (r *TokenRepository) Revoke(ctx context.Context, id uuid.UUID, reason string) error {
	tokenID := pgtype.UUID{
		Bytes: id,
		Valid: true,
	}

	return r.queries.RevokeToken(ctx, db.RevokeTokenParams{
		ID:            tokenID,
		RevokedReason: pgtype.Text{String: reason, Valid: true},
	})
}

func (r *TokenRepository) RevokeByAccessTokenHash(ctx context.Context, accessTokenHash string, reason string) error {
	return r.queries.RevokeTokenByAccessTokenHash(ctx, db.RevokeTokenByAccessTokenHashParams{
		AccessTokenHash: accessTokenHash,
		RevokedReason:   pgtype.Text{String: reason, Valid: true},
	})
}

func (r *TokenRepository) RevokeByAuthorizationCode(ctx context.Context, authorizationCode string, reason string) error {
	return r.queries.RevokeTokensByAuthorizationCode(ctx, db.RevokeTokensByAuthorizationCodeParams{
		AuthorizationCode: pgtype.Text{String: authorizationCode, Valid: true},
		RevokedReason:     pgtype.Text{String: reason, Valid: true},
	})
}

func (r *TokenRepository) UpdateLastUsed(ctx context.Context, id uuid.UUID) error {
	tokenID := pgtype.UUID{
		Bytes: id,
		Valid: true,
	}

	return r.queries.UpdateLastUsedAt(ctx, tokenID)
}

func (r *TokenRepository) mapTokenToDomain(t db.Token) *domain.Token {
	var authCode *string
	if t.AuthorizationCode.Valid {
		authCode = &t.AuthorizationCode.String
	}

	var revokedAt *time.Time
	if t.RevokedAt.Valid {
		revokedAt = &t.RevokedAt.Time
	}

	var revokedReason *string
	if t.RevokedReason.Valid {
		revokedReason = &t.RevokedReason.String
	}

	var lastUsedAt *time.Time
	if t.LastUsedAt.Valid {
		lastUsedAt = &t.LastUsedAt.Time
	}

	return &domain.Token{
		ID:                    t.ID.Bytes,
		AccessTokenHash:       t.AccessTokenHash,
		RefreshTokenHash:      t.RefreshTokenHash,
		AuthorizationCode:     authCode,
		ClientID:              t.ClientID,
		UserID:                t.UserID.Bytes,
		Scopes:                t.Scopes,
		TokenType:             t.TokenType,
		AccessTokenExpiresAt:  t.AccessTokenExpiresAt.Time,
		RefreshTokenExpiresAt: t.RefreshTokenExpiresAt.Time,
		Revoked:               t.Revoked,
		RevokedAt:             revokedAt,
		RevokedReason:         revokedReason,
		CreatedAt:             t.CreatedAt.Time,
		LastUsedAt:            lastUsedAt,
	}
}
