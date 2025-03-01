package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Id        int       `json:"id"`
	Email     string    `json:"email"`
	Password  []byte    `json:"-"` // Password of User stored as a hash
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RefreshTokens struct {
	Id        int             `json:"id"`
	UserId    int             `json:"user_id"`
	TokenHash int             `json:"token_hash"`
	Expiry    jwt.NumericDate `json:"expiry"`
	IsValid   bool            `json:"is_valid"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// Invalid Refresh Tokens
type InvalidTokens struct {
	Id        int       `json:"id"`
	TokenId   int       `json:"token_id"` // Primary ID of refresh token in DB
	TokenHash int       `json:"token_hash"`
	UserId    int       `json:"user_id"` // Primary ID of user in DB
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
