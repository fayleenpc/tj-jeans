package tokenize

import (
	"database/sql"

	"github.com/fayleenpc/tj-jeans/internal/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) CreateBlacklistTokens(t types.Token) error {
	_, err := s.db.Exec(
		"INSERT INTO blacklisted_tokens (token) VALUES (?)", t.Token,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetBlacklistTokenByString(t string) (types.Token, error) {
	// Check if token is blacklisted
	var dbToken types.Token

	err := s.db.QueryRow("SELECT token FROM blacklisted_tokens WHERE token = ?", t).Scan(&dbToken.Token)

	return dbToken, err
}
