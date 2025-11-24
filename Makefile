.PHONY: sqlc sqlc-check

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