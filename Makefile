build:
	@go build -o bin/tj-jeans cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/tj-jeans

templ:
	@templ generate ./platform/web/views

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@, $(MAKECMDGOALS))

migration-up:
	cd "$(CURDIR)/cmd/migrate/" && go run main.go up

migration-down:
	cd "$(CURDIR)/cmd/migrate/" && go run main.go down

gen:
	@protoc \
	--proto_path=internal "internal/types_grpc/types.proto" \
	--go_out=services/common/ --go_opt=paths=source_relative \
	--go-grpc_out=services/common/
	--go-grpc_opt=paths=source_relative