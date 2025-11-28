//go:build integration

package integration

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	testDBName   = "oidc_test"
	testDBUser   = "test_user"
	testDBPass   = "test_password"
	postgresPort = "5432/tcp"
)

// TestDB holds the test database connection and container
type TestDB struct {
	Pool      *pgxpool.Pool
	Container testcontainers.Container
}

// SetupTestDB creates a new PostgreSQL container and returns a connection pool.
// It also runs the schema.sql to create all tables.
func SetupTestDB(t *testing.T) *TestDB {
	t.Helper()
	ctx := context.Background()

	// Get the path to schema.sql
	_, currentFile, _, _ := runtime.Caller(0)
	schemaPath := filepath.Join(filepath.Dir(currentFile), "..", "..", "adapters", "secondary", "postgres", "schema.sql")

	// Create PostgreSQL container
	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(testDBName),
		postgres.WithUsername(testDBUser),
		postgres.WithPassword(testDBPass),
		postgres.WithInitScripts(schemaPath),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	// Get connection string
	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	// Create connection pool
	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		t.Fatalf("failed to parse pool config: %v", err)
	}

	poolConfig.MaxConns = 5
	poolConfig.MinConns = 1

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		t.Fatalf("failed to create pool: %v", err)
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("failed to ping database: %v", err)
	}

	return &TestDB{
		Pool:      pool,
		Container: container,
	}
}

// Teardown closes the pool and terminates the container
func (db *TestDB) Teardown(t *testing.T) {
	t.Helper()
	ctx := context.Background()

	if db.Pool != nil {
		db.Pool.Close()
	}

	if db.Container != nil {
		if err := db.Container.Terminate(ctx); err != nil {
			t.Logf("warning: failed to terminate container: %v", err)
		}
	}
}

// TruncateTables removes all data from tables (respecting FK order)
func (db *TestDB) TruncateTables(t *testing.T) {
	t.Helper()
	ctx := context.Background()

	// Order matters due to foreign key constraints
	tables := []string{
		"tokens",
		"authorization_codes",
		"users",
		"oauth_clients",
	}

	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)
		if _, err := db.Pool.Exec(ctx, query); err != nil {
			t.Fatalf("failed to truncate table %s: %v", table, err)
		}
	}
}

// CleanTable removes all data from a specific table
func (db *TestDB) CleanTable(t *testing.T, tableName string) {
	t.Helper()
	ctx := context.Background()

	query := fmt.Sprintf("TRUNCATE TABLE %s CASCADE", tableName)
	if _, err := db.Pool.Exec(ctx, query); err != nil {
		t.Fatalf("failed to truncate table %s: %v", tableName, err)
	}
}
