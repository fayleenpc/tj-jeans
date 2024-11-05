package cart

import (
	"github.com/fayleenpc/tj-jeans/internal/types"
	pb "github.com/fayleenpc/tj-jeans/services/common/types_grpc"
	"google.golang.org/grpc"
)

type HandlerServer struct {
	pb.UnimplementedOrderServiceServer
	service types.OrderService
}

func NewHandlerServer(grpcServer *grpc.Server, service types.OrderService) {
	handler := &HandlerServer{service: service}
	pb.RegisterOrderServiceServer(grpcServer, handler)
}
