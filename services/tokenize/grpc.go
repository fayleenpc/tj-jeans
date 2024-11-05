package tokenize

import (
	"github.com/fayleenpc/tj-jeans/internal/types"
	pb "github.com/fayleenpc/tj-jeans/services/common/types_grpc"
	"google.golang.org/grpc"
)

type HandlerServer struct {
	pb.UnimplementedTokenServiceServer
	service types.TokenService
}

func NewHandlerServer(grpcServer *grpc.Server, service types.TokenService) {
	handler := &HandlerServer{service: service}
	pb.RegisterTokenServiceServer(grpcServer, handler)
}
