.PHONY: sqlc sqlc-check docker-up docker-down docker-restart docker-logs docker-ps docker-clean docker-rebuild setup dev mocks migrate-up migrate-down setup test test-unit test-integration test-coverage

sqlc:
	cd internal/adapters/secondary/postgres && sqlc generate

sqlc-check:
	cd internal/adapters/secondary/postgres && sqlc compile

migrate-up:
	migrate -path internal/adapters/secondary/postgres/migrations \
		-database "postgresql://user:pass@localhost:5432/oidc?sslmode=disable" up

migrate-down:
	migrate -path internal/adapters/secondary/postgres/migrations \
		-database "postgresql://user:pass@localhost:5432/oidc?sslmode=disable" down 1

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-clean:
	docker-compose down -v

docker-rebuild:
	docker-compose up -d --build --force-recreate

docker-restart:
	docker-compose restart

docker-logs:
	docker-compose logs -f

docker-ps:
	docker-compose ps

dev:
	@go run cmd/server/main.go

mocks:
	@mockery

setup:
	@go install github.com/vektra/mockery/v3@v3.6.1

# Test targets
test: test-unit test-integration

test-unit:
	@echo "Running unit tests..."
	@go test -v -race ./internal/core/services/... ./internal/adapters/...

test-integration:
	@echo "Running integration tests..."
	@go test -v -tags=integration -race ./internal/tests/integration/...

test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"