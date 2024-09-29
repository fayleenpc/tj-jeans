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

func (s *Store) GetBlacklistedTokens() ([]types.Token, error) {
	rows, err := s.db.Query("SELECT * FROM blacklisted_tokens")
	if err != nil {
		return nil, err
	}
	tokens := make([]types.Token, 0)
	for rows.Next() {
		t, err := scanRowIntoBlacklistedTokens(rows)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, *t)
	}
	return tokens, nil
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

func scanRowIntoBlacklistedTokens(rows *sql.Rows) (*types.Token, error) {
	token := new(types.Token)
	err := rows.Scan(
		&token.ID,
		&token.Token,
		&token.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return token, nil
}
