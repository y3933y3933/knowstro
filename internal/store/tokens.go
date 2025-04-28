package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/y3933y3933/knowstro/internal/tokens"
)

type PostgresTokenStore struct {
	db *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{
		db: db,
	}
}

type TokenStore interface {
	Insert(token *tokens.Token) error
	CreateNewToken(userID int, ttl time.Duration, scope string) (*tokens.Token, error)
	DeleteAllTokensForUser(userID int, scope string) error
}

func (t *PostgresTokenStore) Insert(token *tokens.Token) error {
	query := `
		INSERT INTO tokens (hash, user_id, expiry, scope)
		VALUES($1, $2, $3, $4)
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{token.Hash, token.UserID, token.Expiry, token.Scope}

	_, err := t.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (t *PostgresTokenStore) CreateNewToken(userID int, ttl time.Duration, scope string) (*tokens.Token, error) {
	token, err := tokens.GenerateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = t.Insert(token)
	return token, err

}

func (t *PostgresTokenStore) DeleteAllTokensForUser(userID int, scope string) error {

	query := `
	   DELETE FROM tokens
	   WHERE scope = $1 AND user_id = $2
	`

	args := []any{
		scope,
		userID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := t.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
