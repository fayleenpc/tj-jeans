package users

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/fayleenpc/tj-jeans/internal/types"
	pb "github.com/fayleenpc/tj-jeans/services/common/types_grpc"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	rows, err := s.db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	users := make([]*pb.User, 0)
	for rows.Next() {
		userPB, err := scanRowIntoUserPB(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, userPB)
	}
	return &pb.GetUsersResponse{Users: users}, nil
}

func (s *Service) GetUsersByIDs(ctx context.Context, userIDS *pb.GetUsersByIDsRequest) (*pb.GetUsersByIDsResponse, error) {
	placeholders := strings.Repeat(",?", len(userIDS.GetIds())-1)
	query := fmt.Sprintf("SELECT * FROM users WHERE id IN (?%s)", placeholders)
	args := make([]interface{}, len(userIDS.GetIds()))
	for i, v := range userIDS.GetIds() {
		args[i] = v
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	users := []*pb.User{}
	for rows.Next() {
		userPB, err := scanRowIntoUserPB(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, userPB)
	}

	return &pb.GetUsersByIDsResponse{Users: users}, nil
}

func (s *Service) UpdateVerifiedUserByEmail(ctx context.Context, email *pb.UpdateVerifiedUserByEmailRequest) error {
	_, err := s.db.Exec(
		"UPDATE users SET verified = ? WHERE email = ?",
		1, email.GetEmail(),
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetUserByEmail(ctx context.Context, email *pb.GetUserByEmailRequest) (*pb.GetUserByEmailResponse, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email = ?", email.GetEmail())
	if err != nil {
		return nil, err
	}
	userPB := new(pb.User)
	for rows.Next() {
		userPB, err = scanRowIntoUserPB(rows)
		if err != nil {
			return nil, err
		}
	}
	if userPB.GetEmail() != email.GetEmail() {
		return nil, fmt.Errorf("email not found")
	}
	return &pb.GetUserByEmailResponse{User: userPB}, nil
}

func (s *Service) GetUserByID(ctx context.Context, id *pb.GetUserByIDRequest) (*pb.GetUserByIDResponse, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE id = ?", id.GetId())
	if err != nil {
		return nil, err
	}
	userPB := new(pb.User)
	for rows.Next() {
		userPB, err = scanRowIntoUserPB(rows)
		if err != nil {
			return nil, err
		}
	}
	if userPB.GetId() == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return &pb.GetUserByIDResponse{User: userPB}, nil
}

func (s *Service) DeleteUserByID(ctx context.Context, id *pb.DeleteUserByIDRequest) (*pb.DeleteUserByIDResponse, error) {
	res, err := s.db.Exec(
		"DELETE from users WHERE id = ?",
		id.GetId(),
	)
	if err != nil {
		return &pb.DeleteUserByIDResponse{}, err
	}
	deletedID, _ := res.LastInsertId()
	return &pb.DeleteUserByIDResponse{DeletedCount: deletedID}, nil
}

func (s *Service) DeleteUser(ctx context.Context, user *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	res, err := s.db.Exec(
		"DELETE from users WHERE id = ?",
		user.GetUser().GetId(),
	)
	if err != nil {
		return &pb.DeleteUserResponse{}, err
	}
	deletedUser, _ := res.LastInsertId()
	return &pb.DeleteUserResponse{DeletedCount: deletedUser}, nil
}

func (s *Service) UpdateUser(ctx context.Context, user *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	res, err := s.db.Exec(
		"UPDATE users SET firstName = ?, lastName = ?, phoneNumber = ? WHERE id = ?",
		user.GetUser().GetFirstName(), user.GetUser().GetLastName(), user.GetUser().GetPhoneNumber(), user.GetUser().GetId(),
	)

	if err != nil {
		return &pb.UpdateUserResponse{}, err
	}
	updatedUser, _ := res.LastInsertId()
	return &pb.UpdateUserResponse{UpdatedCount: updatedUser}, nil
}

func (s *Service) CreateUser(ctx context.Context, user *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	_, err := s.db.Exec(
		"INSERT INTO users (firstName, lastName, email, password, phoneNumber, address, verified, role) VALUES (?,?,?,?,?,?,?,?)",
		user.GetUser().GetFirstName(), user.GetUser().GetLastName(), user.GetUser().GetEmail(), user.GetUser().GetPassword(), user.GetUser().GetPhoneNumber(), user.GetUser().GetAddress(), user.GetUser().GetVerified(), user.GetUser().GetRole(),
	)
	if err != nil {
		return nil, err
	}
	return &pb.CreateUserResponse{}, nil
}

func scanRowIntoUserPB(rows *sql.Rows) (*pb.User, error) {
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
	userPB := new(pb.User)
	userPB.Address = user.Address
	userPB.CreatedAt = user.CreatedAt.Unix()
	userPB.Email = user.Email
	userPB.FirstName = user.FirstName
	userPB.Id = int32(user.ID)
	userPB.LastName = user.LastName
	userPB.Password = user.Password
	userPB.PhoneNumber = user.PhoneNumber
	userPB.Role = user.Role
	userPB.Verified = user.Verified
	return userPB, nil
}
