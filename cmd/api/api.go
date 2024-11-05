package api

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/fayleenpc/tj-jeans/cmd/docs" // docs is generated by Swag CLI
	"github.com/fayleenpc/tj-jeans/internal/config"
	"github.com/fayleenpc/tj-jeans/internal/monitoring"
	swagger_docs "github.com/fayleenpc/tj-jeans/internal/swaggerdocs"
	"github.com/fayleenpc/tj-jeans/services/cart"
	"github.com/fayleenpc/tj-jeans/services/finance"
	"github.com/fayleenpc/tj-jeans/services/gateway/payment"
	"github.com/fayleenpc/tj-jeans/services/order"
	"github.com/fayleenpc/tj-jeans/services/products"
	"github.com/fayleenpc/tj-jeans/services/tokenize"
	"github.com/fayleenpc/tj-jeans/services/users"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
)

const version = "1.0.0"

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

// test /product post or /product get
// {
// 	"name": "new product12612512",
// 	"description": "new product arrive61251261",
// 	"image" : "new image612612512",
// 	"price": 1000000,
// 	"quantity": 300
// }

// {
// 	"id": 3,
// 	"name": "new jeans",
// 	"description": "new jeans product",
// 	"image": "image new jeans",
// 	"quantity": 100,
// 	"price": 80000
// }

// test /login
// {
//     "email": "me@me.com",
//     "password": "asd"
// }

// test /register
// {
//     "email": "me@me.com",
//     "password": "asd",
//     "firstName": "tiago",
//     "lastName": "user"
// }

// test /cart/checkout
// {
// 	"items": [
// 	  {
// 		"productID": 1,
// 		"quantity": 2
// 	  },
// 	  {
// 		"productID": 2,
// 		"quantity": 3
// 	  }
// 	]
// }

// test /payment/invoices
// {
// 	"payment" : {
// 	  "payment_type" : "alfamart"
// 	},
// 	"customer": {
// 	  "name" : "john",
// 	  "email" : "foo@bar.com",
// 	  "phone_number" : ""
// 	},
// 	"items" : [{
// 	  "name" : "support podcast",
// 	  "category" : "podcast",
// 	  "merchant": "imregi.com",
// 	  "description": "donasi podcast imre",
// 	  "qty": 1,
// 	  "price": 10000,
// 	  "currency": "IDR"
// 	},
// 	{
// 	  "name" : "gk1020h12",
// 	  "category" : "podcast gk1020h12",
// 	  "merchant": "imregi.com gk1020h12",
// 	  "description": "donasi gk1020h12",
// 	  "qty": 1,
// 	  "price": 50000,
// 	  "currency": "IDR"
// 	}]
//   }

func (s *APIServer) Run() error {
	apihost := fmt.Sprintf("%v:%v", config.Envs.PublicHost, config.Envs.Port)

	tracer, closer := monitoring.Jaegar()
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	headersOk := handlers.AllowedHeaders([]string{
		"X-Requested-With",
		"Content-Type, Authorization",
		"Access-Control-Allow-Credentials",
		"Access-Control-Allow-Origin",
		"Access-Control-Allow-Methods",
		"Access-Control-Allow-Headers",
	})
	originsOk := handlers.AllowedOrigins([]string{
		apihost + "/",
		apihost + "/products",
		apihost + "/gallery",
		apihost + "/service",
		apihost + "/api/v1/verify",
		apihost + "/api/v1/login",
		apihost + "/api/v1/refresh",
		apihost + "/api/v1/logout",
		apihost + "/api/v1/register",
		apihost + "/api/v1/payment/*",
		apihost + "/api/v1/cart/checkout",
	})
	methodsOk := handlers.AllowedMethods([]string{
		"GET", "HEAD", "POST", "PUT", "OPTIONS",
	})

	redisStore := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis server address
		Password: "",               // No password set
		DB:       0,                // Use default DB
	})

	// Test the connection
	_, err := redisStore.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis")

	router := mux.NewRouter()

	subrouter := router.PathPrefix("/api/v1").Subrouter()

	// token store
	tokenStore := tokenize.NewStore(s.db)

	usersStore := users.NewStore(s.db)
	usersHandler := users.NewHandler(usersStore, tokenStore, redisStore)
	usersHandler.RegisterRoutes(subrouter)

	productStore := products.NewStore(s.db)
	productHandler := products.NewHandler(productStore, usersStore, tokenStore, redisStore)
	productHandler.RegisterRoutes(subrouter)

	orderStore := order.NewStore(s.db)
	cartHandler := cart.NewHandler(orderStore, productStore, usersStore, tokenStore, redisStore)
	cartHandler.RegisterRoutes(subrouter)

	// payment gateway
	paymentGateway := payment.NewHandler(subrouter, payment.NewServer())
	paymentGateway.RegisterRoutes()

	// tokenize
	tokenizeHandler := tokenize.NewHandler(tokenStore, usersStore, redisStore)
	tokenizeHandler.RegisterRoutes(subrouter)

	// swagger
	// swag := swagger.NewHandler()
	// swag.RegisterRoutes()

	// swagger
	swagDocs := swagger_docs.NewHandler()
	swagDocs.RegisterRoutes(subrouter)

	// mux := http.NewServeMux()
	// handler := NewHandler(c)
	// handler.registerRoutes(mux)
	// log.Println("Starting HTTP server at ", httpAddr)
	// if err := http.ListenAndServe(httpAddr, mux); err != nil {
	// 	log.Fatal("Failed to start http server")
	// }

	financeHandler := finance.NewHandler(orderStore)
	financeHandler.RegisterRoutes(subrouter)

	log.Printf("REST + Json running at : %v\n", s.addr)
	// log.Println("ENVS : ")
	// log.Println(config.Envs)

	// return http.ListenAndServe(s.addr, router)
	creds := handlers.AllowCredentials()

	return http.ListenAndServe(s.addr, handlers.CORS(originsOk, headersOk, methodsOk, creds)(router))
}
