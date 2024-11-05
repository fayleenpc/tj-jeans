package tokenize

import (
	"context"
	"database/sql"

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

func (s *Service) GetBlacklistedTokens(ctx context.Context, req *pb.GetBlacklistedTokensRequest) (*pb.GetBlacklistedTokensResponse, error) {
	rows, err := s.db.Query("SELECT * FROM blacklisted_tokens")
	if err != nil {
		return nil, err
	}
	tokensPB := make([]*pb.Token, 0)
	for rows.Next() {
		t, err := scanRowIntoBlacklistedTokensPB(rows)
		if err != nil {
			return nil, err
		}
		tokensPB = append(tokensPB, t)
	}
	return &pb.GetBlacklistedTokensResponse{Tokens: tokensPB}, nil
}

func (s *Service) CreateBlacklistTokens(ctx context.Context, t *pb.CreateBlacklistTokenRequest) (*pb.CreateBlacklistTokenResponse, error) {
	return nil, nil
}

func (s *Service) GetBlacklistTokenByString(ctx context.Context, t *pb.GetBlacklistTokenByStringRequest) (*pb.GetBlacklistTokenByStringResponse, error) {
	return nil, nil
}

func scanRowIntoBlacklistedTokensPB(rows *sql.Rows) (*pb.Token, error) {
	token := new(types.Token)
	err := rows.Scan(
		&token.ID,
		&token.Token,
		&token.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	tokenPB := new(pb.Token)
	tokenPB.Id = int32(token.ID)
	tokenPB.Token = token.Token
	tokenPB.CreatedAt = token.CreatedAt.Unix()
	return tokenPB, nil
}
