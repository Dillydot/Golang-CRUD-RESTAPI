DB_URL = "postgres://test-user:a12345@localhost:5432/test-db?sslmode=disable"

fmt:
	go fmt ./...

# Generate Go mocks
mocks:
	mockgen -source ./internal/app/service.go -destination ./internal/app/service_mocks.go -package app
	mockgen -source ./internal/infra/http/router.go -destination ./internal/infra/http/router_mocks.go -package http

# https://github.com/golang-migrate/migrate
db-migrate:
	migrate -source "file://db/migrations" -database $(DB_URL) up
