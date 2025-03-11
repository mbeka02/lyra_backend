package database

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type terminateFunc func(context.Context, ...testcontainers.TerminateOption) error

func setupPostgres(ctx context.Context) (testcontainers.Container, string, error) {
	var (
		dbName   = "testdb"
		dbPwd    = "password"
		dbUser   = "user"
		dbSchema = "public"
	)

	pgContainer, err := postgres.Run(
		ctx,
		"postgres:latest",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPwd),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, "", err
	}

	dbHost, err := pgContainer.Host(ctx)
	if err != nil {
		return pgContainer, "", err
	}

	dbPort, err := pgContainer.MappedPort(ctx, "5432/tcp")
	if err != nil {
		return pgContainer, "", err
	}

	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s",
		dbUser, dbPwd, dbHost, dbPort.Port(), dbName, dbSchema)

	return pgContainer, connString, nil
}

func mustStartPostgresContainer() (terminateFunc, error) {
	pgContainer, connString, err := setupPostgres(context.Background())
	if err != nil {
		log.Fatalf("unable to setup postgres")
	}
	SetConnectionString(connString)
	return pgContainer.Terminate, err
}

func TestMain(m *testing.M) {
	teardown, err := mustStartPostgresContainer()
	if err != nil {
		log.Fatalf("could not start postgres container: %v", err)
	}

	m.Run()

	if teardown != nil && teardown(context.Background()) != nil {
		log.Fatalf("could not teardown postgres container: %v", err)
	}
}

func TestNewStore(t *testing.T) {
	srv := NewStore()
	if srv == nil {
		t.Fatal("New() returned nil")
	}
}

func TestStoreHealth(t *testing.T) {
	srv := NewStore()

	stats := srv.Health()

	if stats["status"] != "up" {
		t.Fatalf("expected status to be up, got %s", stats["status"])
	}

	if _, ok := stats["error"]; ok {
		t.Fatalf("expected error not to be present")
	}

	if stats["message"] != "It's healthy" {
		t.Fatalf("expected message to be 'It's healthy', got %s", stats["message"])
	}
}

func TestCloseStore(t *testing.T) {
	srv := NewStore()

	if srv.Close() != nil {
		t.Fatalf("expected Close() to return nil")
	}
}
