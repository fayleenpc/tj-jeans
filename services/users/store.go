package users

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/fayleenpc/tj-jeans/internal/types"
)

// signature
// GetUsers() ([]types.User, error)
// GetUsersByIDs(userIDS []int) ([]types.User, error)
// UpdateVerifiedUserByEmail(email string) error
// GetUserByEmail(email string) (*types.User, error)
// GetUserByID(id int) (*types.User, error)
// DeleteUserByID(id int) (int64, error)
// DeleteUser(user types.User)
// UpdateUser(user types.User) (int64, error)
// CreateUser(user types.User) error

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) GetUsers() ([]types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	users := make([]types.User, 0)
	for rows.Next() {
		user, err := scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}
	return users, nil
}

func (s *Store) GetUsersByIDs(userIDS []int) ([]types.User, error) {
	placeholders := strings.Repeat(",?", len(userIDS)-1)
	query := fmt.Sprintf("SELECT * FROM users WHERE id IN (?%s)", placeholders)

	args := make([]interface{}, len(userIDS))
	for i, v := range userIDS {
		args[i] = v
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	users := []types.User{}
	for rows.Next() {
		user, err := scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, *user)
	}

	return users, nil
}

func (s *Store) UpdateVerifiedUserByEmail(email string) error {
	_, err := s.db.Exec(
		"UPDATE users SET verified = ? WHERE email = ?",
		1, email,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}
	user := new(types.User)
	for rows.Next() {
		user, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}
	if user.Email != email {
		return nil, fmt.Errorf("email not found")
	}
	return user, nil
}

func (s *Store) GetUserByID(id int) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	user := new(types.User)
	for rows.Next() {
		user, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}
	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (s *Store) DeleteUserByID(id int) (int64, error) {
	res, err := s.db.Exec(
		"DELETE from users WHERE id = ?",
		id,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) DeleteUser(user types.User) (int64, error) {
	res, err := s.db.Exec(
		"DELETE from users WHERE id = ?",
		user.ID,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) UpdateUser(user types.User) (int64, error) {
	res, err := s.db.Exec(
		"UPDATE users SET firstName = ?, lastName = ?, phoneNumber = ? WHERE id = ?",
		user.FirstName, user.LastName, user.PhoneNumber, user.ID,
	)

	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (s *Store) CreateUser(user types.User) error {
	_, err := s.db.Exec(
		"INSERT INTO users (firstName, lastName, email, password, phoneNumber, address, verified, role) VALUES (?,?,?,?,?,?,?,?)",
		user.FirstName, user.LastName, user.Email, user.Password, user.PhoneNumber, user.Address, user.Verified, user.Role,
	)
	if err != nil {
		return err
	}
	return nil
}

func scanRowIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)
	err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.PhoneNumber,
		&user.Address,
		&user.Verified,
		&user.Role,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}
