package repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/mbeka02/lyra_backend/internal/database"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var store *database.Store

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

func TestMain(m *testing.M) {
	code := runTests(m)
	os.Exit(code)
}

func runTests(m *testing.M) int {
	// setup postgres test container
	ctx := context.Background()
	pgContainer, connString, err := setupPostgres(ctx)
	if err != nil {
		log.Fatalf("unable to setup the test container:%v", err)
	}
	// set connection string
	database.SetConnectionString(connString)

	// initialize the store
	store = database.NewStore()
	// run migrations
	if err := runMigrations(store); err != nil {
		log.Fatal(err)
	}
	// cleanup
	defer func() {
		if err := store.Close(); err != nil {
			log.Printf("failed to close the db:%v", err)
		}
		if err := pgContainer.Terminate(ctx); err != nil {
			log.Printf("could not terminate the test container:%v", err)
		}
	}()
	// run tests and return a code
	return m.Run()
}

func runMigrations(store *database.Store) error {
	// get the DB instance
	dbInstance := store.DB()
	// setup goose

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("unable to setup dialect for the goose migration tool:%v", err)
	}
	err := goose.Up(dbInstance, "../../../sql/schema")
	if err != nil {
		return fmt.Errorf("unable to  execute db migration: %v", err)
	}
	return nil
}
