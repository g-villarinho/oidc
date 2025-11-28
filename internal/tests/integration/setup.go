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
	"github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	redisContainer "github.com/testcontainers/testcontainers-go/modules/redis"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	testDBName   = "oidc_test"
	testDBUser   = "test_user"
	testDBPass   = "test_password"
	postgresPort = "5432/tcp"
	redisPort    = "6379/tcp"
)

type TestDB struct {
	Pool      *pgxpool.Pool
	Container testcontainers.Container
}

type TestRedis struct {
	Client    *redis.Client
	Container testcontainers.Container
}

type TestEnv struct {
	DB    *TestDB
	Redis *TestRedis
}

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

// SetupTestRedis creates and starts a Redis test container
func SetupTestRedis(t *testing.T) *TestRedis {
	t.Helper()
	ctx := context.Background()

	// Create Redis container
	container, err := redisContainer.Run(ctx,
		"redis:7-alpine",
		testcontainers.WithWaitStrategy(
			wait.ForLog("Ready to accept connections").
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("failed to start redis container: %v", err)
	}

	// Get connection string
	connStr, err := container.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("failed to get redis connection string: %v", err)
	}

	// Create Redis client
	opts, err := redis.ParseURL(connStr)
	if err != nil {
		t.Fatalf("failed to parse redis URL: %v", err)
	}

	client := redis.NewClient(opts)

	// Verify connection
	if err := client.Ping(ctx).Err(); err != nil {
		t.Fatalf("failed to ping redis: %v", err)
	}

	return &TestRedis{
		Client:    client,
		Container: container,
	}
}

// Teardown closes the client and terminates the container
func (r *TestRedis) Teardown(t *testing.T) {
	t.Helper()
	ctx := context.Background()

	if r.Client != nil {
		if err := r.Client.Close(); err != nil {
			t.Logf("warning: failed to close redis client: %v", err)
		}
	}

	if r.Container != nil {
		if err := r.Container.Terminate(ctx); err != nil {
			t.Logf("warning: failed to terminate redis container: %v", err)
		}
	}
}

// FlushAll removes all data from Redis
func (r *TestRedis) FlushAll(t *testing.T) {
	t.Helper()
	ctx := context.Background()

	if err := r.Client.FlushAll(ctx).Err(); err != nil {
		t.Fatalf("failed to flush redis: %v", err)
	}
}

// SetupTestEnv creates a complete test environment with PostgreSQL and Redis
func SetupTestEnv(t *testing.T) *TestEnv {
	t.Helper()

	db := SetupTestDB(t)
	redis := SetupTestRedis(t)

	return &TestEnv{
		DB:    db,
		Redis: redis,
	}
}

// Teardown closes all resources and terminates containers
func (env *TestEnv) Teardown(t *testing.T) {
	t.Helper()

	if env.DB != nil {
		env.DB.Teardown(t)
	}

	if env.Redis != nil {
		env.Redis.Teardown(t)
	}
}

// Reset cleans all data from databases
func (env *TestEnv) Reset(t *testing.T) {
	t.Helper()

	if env.DB != nil {
		env.DB.TruncateTables(t)
	}

	if env.Redis != nil {
		env.Redis.FlushAll(t)
	}
}
