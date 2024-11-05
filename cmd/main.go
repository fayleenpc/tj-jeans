package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/fayleenpc/tj-jeans/cmd/api"
	"github.com/fayleenpc/tj-jeans/cmd/api_proto"
	"github.com/fayleenpc/tj-jeans/internal/config"
	"github.com/fayleenpc/tj-jeans/internal/db"
	"github.com/fayleenpc/tj-jeans/platform/web"
	"github.com/fayleenpc/tj-jeans/services/tokenize"
	"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

var (
	grpcAddr = ":2000"
)

func main() {

	db, err := db.NewMySQLStorage(mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	})
	if err != nil {
		log.Fatal(err)
	}
	initStorage(db)

	// gRPC API
	// grpcApiServer := api_grpc.NewApiServerGRPC(":8082", grpc.NewServer(), db)

	// go grpcApiServer.Run()

	// REST PROTOBUF API
	restApiProtobuf := api_proto.NewApiProtobufServer(":8082", db)

	go restApiProtobuf.Run()

	// REST API
	restApiServer := api.NewAPIServer(":8081", db)

	go restApiServer.Run()

	startWeb(db)
}

func startWeb(db *sql.DB) {
	tokenStore := tokenize.NewStore(db)
	router := mux.NewRouter()

	// serve files in static folder
	router.PathPrefix("/platform/web/static/").Handler(http.StripPrefix("/platform/web/static/", http.FileServer(http.Dir("platform/web/static"))))
	router.PathPrefix("/platform/web/static/images/").Handler(http.StripPrefix("/platform/web/static/images/", http.FileServer(http.Dir("platform/web/static/images"))))

	// servefiles in static_admin folder
	router.PathPrefix("/platform/web/static_admin/").Handler(http.StripPrefix("/platform/web/static_admin/", http.FileServer(http.Dir("platform/web/static_admin"))))
	router.PathPrefix("/platform/web/static_admin/images/").Handler(http.StripPrefix("/platform/web/static_admin/images/", http.FileServer(http.Dir("platform/web/static_admin/images"))))

	// web
	web := web.NewHandler(tokenStore)
	web.RegisterRoutes(router)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB: successfully connected!")
}
