package repository

import (
	"auth-service/internal/model"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshTokenRepository struct {
	pool *pgxpool.Pool
}

func NewRefreshTokenRepository(pool *pgxpool.Pool) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		pool: pool,
	}
}

func (r *RefreshTokenRepository) CreateRefreshToken(newRefreshToken model.RefreshToken) (model.RefreshToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at) 
		VALUES($1, $2, $3)
		RETURNING id, user_id, token_hash, expires_at, created_at, revoked_at;
	`

	var refreshToken model.RefreshToken

	err := r.pool.QueryRow(ctx, query,
		newRefreshToken.UserID,
		newRefreshToken.TokenHash,
		newRefreshToken.ExpiresAt,
	).Scan(
		&refreshToken.ID,
		&refreshToken.UserID,
		&refreshToken.TokenHash,
		&refreshToken.ExpiresAt,
		&refreshToken.CreatedAt,
		&refreshToken.RevokedAt,
	)
	if err != nil {
		return model.RefreshToken{}, err
	}

	return refreshToken, nil
}

func (r *RefreshTokenRepository) GetRefreshTokenByHash(tokenHash string) (model.RefreshToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, user_id, token_hash, expires_at, created_at, revoked_at
		FROM refresh_tokens
		WHERE token_hash = $1;
	`

	var refreshToken model.RefreshToken

	err := r.pool.QueryRow(ctx, query, tokenHash).Scan(
		&refreshToken.ID,
		&refreshToken.UserID,
		&refreshToken.TokenHash,
		&refreshToken.ExpiresAt,
		&refreshToken.CreatedAt,
		&refreshToken.RevokedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.RefreshToken{}, ErrRefreshTokenNotFound
	}

	if err != nil {
		return model.RefreshToken{}, err
	}

	return refreshToken, nil
}

func (r *RefreshTokenRepository) RevokeRefreshToken(tokenHash string) (model.RefreshToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	revokedAt := time.Now()

	query := `
		UPDATE refresh_tokens
		SET revoked_at = $2
		WHERE token_hash = $1 AND revoked_at IS NULL
		RETURNING id, user_id, token_hash, expires_at, created_at, revoked_at;
	`
	var refreshToken model.RefreshToken
	err := r.pool.QueryRow(ctx, query, tokenHash, revokedAt).Scan(
		&refreshToken.ID,
		&refreshToken.UserID,
		&refreshToken.TokenHash,
		&refreshToken.ExpiresAt,
		&refreshToken.CreatedAt,
		&refreshToken.RevokedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.RefreshToken{}, ErrRefreshTokenNotFound
	}

	if err != nil {
		return model.RefreshToken{}, err
	}

	return refreshToken, nil
}

func (r *RefreshTokenRepository) RevokeAllRefreshTokens(userID uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	revokedAt := time.Now()

	query := `
		UPDATE refresh_tokens
		SET revoked_at = $2
		WHERE user_id = $1 AND revoked_at IS NULL;
	`

	_, err := r.pool.Exec(ctx, query, userID, revokedAt)
	if err != nil {
		return err
	}

	return nil
}
