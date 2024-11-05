package products

import (
	"github.com/fayleenpc/tj-jeans/internal/types"
	pb "github.com/fayleenpc/tj-jeans/services/common/types_grpc"
	"google.golang.org/grpc"
)

type HandlerServer struct {
	pb.UnimplementedProductServiceServer
	service types.ProductService
}

func NewHandlerServer(grpcServer *grpc.Server, service types.ProductService) {
	handler := &HandlerServer{service: service}

	pb.RegisterProductServiceServer(grpcServer, handler)
}
