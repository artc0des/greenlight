package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"time"

	"greenlight.artc0des.com/internal/validator"
)

const (
	ScopeActivation = "activation"
)

type Token struct {
	Plaintext string
	Hash      []byte
	UserID    int64
	Expiry    time.Time
	Scope     string
}

func generateToken(userID int64, ttl time.Duration, scope string) *Token {
	token := &Token{
		Plaintext: rand.Text(),
		UserID:    userID,
		Expiry:    time.Now().Add(ttl),
		Scope:     scope,
	}

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 characters long")
}

type TokenModel struct {
	DB *sql.DB
}

func (m *TokenModel) New(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token := generateToken(userID, ttl, scope)

	err := m.Insert(token)
	return token, err
}

func (m *TokenModel) Insert(token *Token) error {
	query := `insert into tokens (hash, user_id, expiry, scope)
				values ($1, $2, $3, $4)`

	args := []any{token.Hash, token.UserID, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}
