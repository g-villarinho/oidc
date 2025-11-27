.PHONY: sqlc sqlc-check docker-up docker-down docker-restart docker-logs docker-ps docker-clean docker-rebuild setup dev

sqlc:
	sqlc generate

sqlc-check:
	sqlc compile

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