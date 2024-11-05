package users

import (
	"github.com/fayleenpc/tj-jeans/internal/types"
	pb "github.com/fayleenpc/tj-jeans/services/common/types_grpc"
	"google.golang.org/grpc"
)

type HandlerServer struct {
	pb.UnimplementedUserServiceServer
	service types.UserService
}

func NewHandlerServer(grpcServer *grpc.Server, service types.UserService) {
	handler := &HandlerServer{service: service}
	pb.RegisterUserServiceServer(grpcServer, handler)
}
