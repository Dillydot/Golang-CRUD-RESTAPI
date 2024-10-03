
fmt:
	go fmt ./...

# Generate Go mocks
# https://github.com/uber-go/mock
mocks:
	mockgen -source ./internal/app/service.go -destination ./internal/app/service_mocks.go -package app
	mockgen -source ./internal/infra/http/router.go -destination ./internal/infra/http/router_mocks.go -package http

run:
	docker-compose -f docker-compose.yml up -d --build
