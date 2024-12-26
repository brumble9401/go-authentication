build:
	@go build -o bin/golang-authentication cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/golang-authentication

migrate-up:
	@go run cmd/migrate/main.go