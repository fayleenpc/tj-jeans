package api_proto

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/fayleenpc/tj-jeans/internal/monitoring"
	"github.com/fayleenpc/tj-jeans/services/cart"
	"github.com/fayleenpc/tj-jeans/services/order"
	"github.com/fayleenpc/tj-jeans/services/products"
	"github.com/fayleenpc/tj-jeans/services/tokenize"
	"github.com/fayleenpc/tj-jeans/services/users"
	"github.com/opentracing/opentracing-go"
)

type ApiProtobufServer struct {
	addr string
	db   *sql.DB
}

func NewApiProtobufServer(addr string, db *sql.DB) *ApiProtobufServer {
	return &ApiProtobufServer{
		addr: addr,
		db:   db,
	}
}

func (s *ApiProtobufServer) Run() error {

	tracer, closer := monitoring.Jaegar()
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	usersService := users.NewService(s.db)
	tokenizeService := tokenize.NewService(s.db)
	productsService := products.NewService(s.db)
	ordersService := order.NewService(s.db)

	mux := http.NewServeMux()

	// orders service
	orderHandlerHTTP := cart.NewHandlerHTTP(ordersService)
	orderHandlerHTTP.RegisterRoutes(mux)

	// products servuce
	productsHandlerHTTP := products.NewHandlerHTTP(productsService)
	productsHandlerHTTP.RegisterRoutes(mux)

	// tokenize service
	tokenizeHandlerHTTP := tokenize.NewHandlerHTTP(tokenizeService)
	tokenizeHandlerHTTP.RegisterRoutes(mux)

	// users service
	usersHandlerHTTP := users.NewHandlerHTTP(usersService)
	usersHandlerHTTP.RegisterRoutes(mux)

	log.Printf("REST + Protobuf running at : %v\n", s.addr)

	return http.ListenAndServe(s.addr, mux)
}
