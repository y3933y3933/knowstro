package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgconn"
	"github.com/y3933y3933/knowstro/internal/tokens"
	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plainText *string
	hash      []byte
}

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"username"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

var AnonymousUser = &User{}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{
		db: db,
	}
}

type UserStore interface {
	CreateUser(*User) error
	UpdateUser(*User) error
	GetUserByName(username string) (*User, error)
	GetForToken(tokenScope, tokenPlaintext string) (*User, error)
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.plainText = &plaintextPassword
	p.hash = hash
	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))

	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func (s *PostgresUserStore) CreateUser(user *User) error {
	query := `
		INSERT INTO users (name, email, password_hash, activated)
		VALUES($1, $2, $3, $4)
		RETURNING id, created_at, version
	`
	args := []any{user.Name, user.Email, user.Password.hash, user.Activated}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == UniqueViolationErr {
			if pgErr.ConstraintName == "users_email_key" {
				return ErrDuplicateEmail
			}
			if pgErr.ConstraintName == "users_name_unique" {
				return ErrDuplicateUserName
			}
		}

		return err
	}
	return nil
}

func (s *PostgresUserStore) GetUserByName(username string) (*User, error) {
	user := &User{
		Password: password{},
	}

	query := `
		SELECT id, username, email, password_hash, activated, created_at, version
		FROM users
		WHERE name = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.CreatedAt,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}

	}

	return user, nil
}

func (s *PostgresUserStore) GetForToken(tokenScope, tokenPlaintext string) (*User, error) {
	hash := tokens.HashTokenPlainText(tokenPlaintext)

	query := `
		SELECT users.id, users.created_at, users.name, users.email, users.password_hash,users.activated
		FROM users
		INNER JOIN tokens
		ON users.id = tokens.user_id
		WHERE tokens.hash = $1
		AND tokens.scope = $2
		AND tokens.expiry > $3
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User

	args := []any{
		hash,
		tokenScope,
		time.Now(),
	}

	err := s.db.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil

}

func (s *PostgresUserStore) UpdateUser(user *User) error {
	query := `
		UPDATE users
		SET name = $1, email = $2, password_hash = $3, activated = $4, version = version + 1
		WHERE id = $5 AND version = $6
		RETURNING version
	`

	args := []any{
		user.Name,
		user.Email,
		user.Password.hash,
		user.Activated,
		user.ID,
		user.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == UniqueViolationErr {
			if pgErr.ConstraintName == "users_email_key" {
				return ErrDuplicateEmail
			}
			if pgErr.ConstraintName == "users_name_unique" {
				return ErrDuplicateUserName
			}
		}

		return err
	}

	return nil
}
