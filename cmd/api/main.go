package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/mbeka02/lyra_backend/config"
	"github.com/mbeka02/lyra_backend/internal/auth"
	"github.com/mbeka02/lyra_backend/internal/imgstore"
	"github.com/mbeka02/lyra_backend/internal/server"
)

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func setupServer() (*http.Server, error) {
	conf, err := config.LoadConfig(".")
	if err != nil {
		return nil, fmt.Errorf("unable to load config: %v", err)
	}

	maker, err := auth.NewJWTMaker(conf.SYMMETRIC_KEY)
	if err != nil {
		return nil, fmt.Errorf("unable to setup the auth token maker:%v", err)
	}

	ImageFileStorage, err := imgstore.NewGCStorage(conf.GCLOUD_PROJECT_ID, conf.GCLOUD_BUCKET_NAME)
	if err != nil {
		return nil, fmt.Errorf("unable to setup cloud storage:%v", err)
	}
	// Configutation Options
	opts := server.ConfigOptions{
		Port:                conf.PORT,
		AccessTokenDuration: conf.ACCESS_TOKEN_DURATION,
	}
	server := server.NewServer(opts, maker, ImageFileStorage)
	return server, nil
}

func main() {
	server, err := setupServer()
	if err != nil {
		log.Fatalf("fatal error , the server setup process failed : %v", err)
	}
	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(server, done)

	log.Println("the server is listening on port" + server.Addr)
	// mailer.SendEmail()
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")
}
