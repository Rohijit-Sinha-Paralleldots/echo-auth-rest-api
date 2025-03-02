package models

import (
	"database/sql"
	"math/rand"
	"time"

	"github.com/Rohijit-Sinha-Paralleldots/echo-auth-rest-api/storage"
	"github.com/golang-jwt/jwt/v5"
)

type UserSchema struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Database Schema of `users` table
type User struct {
	Id        int       `json:"id"`
	Email     string    `json:"email"`
	Password  []byte    `json:"-"` // Password of User stored as a hash
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (*User) Create(user UserSchema) (User, error) {
	db := storage.GetDB()
	resUser := User{}
	sqlStmt := `insert into users (email, password, created_at, updated_at) values ($1, $2, $3, $4)\
	returning id, email, created_at, modified_at`
	curTime := time.Now().UTC()
	err := db.QueryRow(sqlStmt, user.Email, user.Password, curTime, curTime).Scan(&resUser)
	if err != nil {
		return resUser, err
	}
	return resUser, nil
}

// Database Schema of `refresh_tokens` table
type RefreshToken struct {
	Id        int       `json:"id"`
	UserId    int       `json:"user_id"`
	TokenHash int       `json:"token_hash"`
	Expiry    time.Time `json:"expiry"`
	IsValid   int       `json:"is_valid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (*RefreshToken) Create(user User) (RefreshToken, error) {
	db := storage.GetDB()
	tokenHash := rand.Int63()
	curTime := time.Now().UTC()
	expiry := jwt.NewNumericDate(curTime.Add(time.Hour * 24 * 15))
	resToken := RefreshToken{}
	sqlStmt := `insert into refresh_tokens (user_id, token_hash, expiry, is_valid,created_at, updated_at) \
	values ($1, $2, $3, $4, $5, $6) returning id, user_id, token_hash, expiry, is_valid, created_at, updated_at`
	err := db.QueryRow(sqlStmt, user.Id, tokenHash, expiry.Time, 1, curTime, curTime).Scan(&resToken)
	if err != nil {
		return resToken, err
	}
	return resToken, nil
}

func (token *RefreshToken) Invalidate(tx *sql.Tx) error {
	curTime := time.Now()
	sqlStmt1 := `update refresh_tokens where id = $1 set is_valid = 0, updated_at = $2`
	_, err := tx.Exec(sqlStmt1, token.Id, curTime.UTC())
	if err != nil {
		return err
	}
	return nil
}

// Database schema of `invalid_refresh_tokens` tables
type InvalidToken struct {
	Id        int       `json:"id"`
	TokenId   int       `json:"token_id"` // Primary ID of refresh token in DB
	TokenHash int       `json:"token_hash"`
	UserId    int       `json:"user_id"` // Primary ID of user in DB
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (*InvalidToken) Create(token RefreshToken, user User) (InvalidToken, error) {
	resToken := InvalidToken{}
	db := storage.GetDB()
	curTime := time.Now()
	dbTx, err := db.Begin()
	if err != nil {
		return resToken, err
	}
	err = token.Invalidate(dbTx)
	if err != nil {
		dbTx.Rollback()
		return resToken, err
	}
	sqlStmt2 := `insert into invalid_refresh_tokens
				(token_id, token_hash, user_id, created_at, updated_at)
				values ($1, $2, $3, $4, $5)
				returning id, token_id, token_hash, user_id, created_at, updated_at`
	err = dbTx.QueryRow(sqlStmt2, token.Id, token.TokenHash, token.UserId, curTime.UTC(), curTime.UTC()).Scan(&resToken)
	if err != nil {
		dbTx.Rollback()
		return resToken, err
	}
	err = dbTx.Commit()
	if err != nil {
		return resToken, err
	}
	return resToken, nil
}
