include .env
# Build the application
all: build

DB_URL=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5432/${POSTGRES_DB}?sslmode=disable

build:
	@echo "Building..."


	@go build -o main cmd/api/main.go

run:
	@go run cmd/api/main.go

generate:
	@sqlc generate

up:
	~/go/bin/goose -dir ./database/schema postgres "$(DB_URL)" up

status:
	goose -dir ./database/schema postgres "$(DB_URL)" status

down:
	~/go/bin/goose -dir ./database/schema postgres "$(DB_URL)" down
	
fixtures:
	goose -dir ./database/schema postgres "$(DB_URL)" up