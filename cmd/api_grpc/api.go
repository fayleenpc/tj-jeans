package api_grpc

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	"github.com/fayleenpc/tj-jeans/internal/monitoring"
	"github.com/fayleenpc/tj-jeans/services/cart"
	pb "github.com/fayleenpc/tj-jeans/services/common/types_grpc"
	"github.com/fayleenpc/tj-jeans/services/order"
	"github.com/fayleenpc/tj-jeans/services/products"
	"github.com/fayleenpc/tj-jeans/services/tokenize"
	"github.com/fayleenpc/tj-jeans/services/users"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

type ApiServerGRPC struct {
	srv  *grpc.Server
	db   *sql.DB
	addr string
}

func NewApiServerGRPC(addr string, grpcServer *grpc.Server, db *sql.DB) *ApiServerGRPC {
	return &ApiServerGRPC{
		srv:  grpcServer,
		db:   db,
		addr: addr,
	}
}

func (s *ApiServerGRPC) Run() {
	s.RunServer()
	defer s.RunClient()
}

func (s *ApiServerGRPC) RunServer() error {
	tracer, closer := monitoring.Jaegar()
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatal("Failed to listen : ", err)
	}
	defer lis.Close()

	usersService := users.NewService(s.db)
	tokenService := tokenize.NewService(s.db)
	productsService := products.NewService(s.db)
	ordersService := order.NewService(s.db)

	users.NewHandlerServer(s.srv, usersService)
	tokenize.NewHandlerServer(s.srv, tokenService)
	products.NewHandlerServer(s.srv, productsService)
	cart.NewHandlerServer(s.srv, ordersService)
	log.Printf("gRPC server is running at : %v\n", s.addr)

	return s.srv.Serve(lis)
}

func (s *ApiServerGRPC) RunClient() error {

	connProducts, err := grpc.NewClient(s.addr)
	if err != nil {
		log.Fatal("Failed to listen on connection [users]")
	}
	defer connProducts.Close()
	productsConnection := pb.NewProductServiceClient(connProducts)
	httpHandler := products.NewHandlerClient(productsConnection)
	response, err := httpHandler.Client.GetProducts(context.Background(), &pb.GetProductsRequest{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("response result [products grpc] = ", response.GetProducts())
	// Create Request
	// products, err := httpHandler.GetProducts(context.Background(), &pb.GetProductsRequest{})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println("Got products from gRPC : ", products)

	// Create Basic Request
	// products, err := productsConnection.GetProducts(context.Background(), &pb.GetProductsRequest{})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println("Got products from gRPC : ", products)

	return err
}
