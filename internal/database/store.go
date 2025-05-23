package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

//
// // Store represents a service that interacts with a database.
// type Store interface {
// 	// Health returns a map of health status information.
// 	// The keys and values in the map are service-specific.
// 	Health() map[string]string
//
// 	// Close terminates the database connection.
// 	// It returns an error if the connection cannot be closed.
// 	Close() error
// }

type Store struct {
	db *sql.DB
	*Queries
}

var (
	connStr    = os.Getenv("DB_CONNECTION_STRING")
	port       = os.Getenv("DB_PORT")
	dbInstance *Store
)

func NewStore() *Store {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}
	conn, err := sql.Open("pgx", connStr)
	// conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}

	// Configure connection pool settings
	// conn.SetMaxOpenConns(50) // Maximum number of open connections
	// conn.SetMaxIdleConns(10) // Maximum number of idle connections
	// conn.SetConnMaxLifetime(5 * time.Minute) // Lifetime of each connection

	dbInstance = &Store{
		db:      conn,
		Queries: New(conn),
	}
	return dbInstance
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *Store) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *Store) Close() error {
	log.Println("Disconnected from the db")
	return s.db.Close()
}

// overrides the default connection string
func SetConnectionString(str string) {
	connStr = str
}

// expose the db instance
func (s *Store) DB() *sql.DB {
	return s.db
}

// executes queries  within a db transaction
func (s *Store) ExecTx(ctx context.Context, fn func(*Queries) error) error {
	// get a tx for making transaction requests
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	// exec the callback fn and return an error if it fails
	err = fn(q)
	// rollback the transaction in case of failure
	if err != nil {
		// if the rollback also fails return both errors
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf(" transaction error:%v,rollback error:%v", err, rbErr)
		}
		return err
	}
	// commit the transaction and return its error if it occurs
	return tx.Commit()
}
