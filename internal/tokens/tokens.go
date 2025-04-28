package tokens

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"
)

const (
	ScopeActivation = "activation"
)

type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    int       `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

func GenerateToken(userID int, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	emptyBytes := make([]byte, 32)

	_, err := rand.Read(emptyBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(emptyBytes)
	token.Hash = HashTokenPlainText(token.Plaintext)
	return token, nil
}

func HashTokenPlainText(tokenPlaintext string) []byte {
	hash := sha256.Sum256([]byte(tokenPlaintext))
	return hash[:]
}
