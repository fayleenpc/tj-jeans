build:
	@go build -o bin/ecom cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/ecom

templ:
	@templ generate ./views

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@, $(MAKECMDGOALS))

migration-up:
	cd "$(CURDIR)/cmd/migrate/" && go run main.go up

migration-down:
	cd "$(CURDIR)/cmd/migrate/" && go run main.go down