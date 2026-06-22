package model

import "time"

type RefreshToken struct {
	ID        uint64     `json:"id"`
	UserID    uint64     `json:"user_id"`
	TokenHash string     `json:"token_hash"`
	ExpiresAt time.Time  `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
	RevokedAt *time.Time `json:"revoked_at"`
}
