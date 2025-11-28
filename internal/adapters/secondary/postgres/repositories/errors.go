package repositories

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func isUniqueViolation(err error) bool {
	// pgx retorna um código específico para unique violation
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505" // unique_violation
	}
	return false
}

func isNotFound(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
