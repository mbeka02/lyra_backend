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
	"github.com/mbeka02/lyra_backend/internal/fhir"
	"github.com/mbeka02/lyra_backend/internal/objstore"
	"github.com/mbeka02/lyra_backend/internal/payment"
	"github.com/mbeka02/lyra_backend/internal/server"
	"github.com/mbeka02/lyra_backend/internal/streamsdk"
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
	// auth setup
	maker, err := auth.NewJWTMaker(conf.SYMMETRIC_KEY)
	if err != nil {
		return nil, fmt.Errorf("unable to setup the auth token maker:%v", err)
	}
	// cloud storage setup for images
	imgStorage, err := objstore.NewGCStorage(conf.GCLOUD_PROJECT_ID, conf.GCLOUD_IMAGE_BUCKET)
	if err != nil {
		return nil, fmt.Errorf("unable to setup image cloud storage:%v", err)
	}
	// cloud storage setup for files
	fileStorage, err := objstore.NewGCStorage(conf.GCLOUD_PROJECT_ID, conf.GCLOUD_PATIENT_RECORD_BUCKET)
	if err != nil {
		return nil, fmt.Errorf("unable to setup file cloud storage:%v", err)
	}
	// setup the fhir client
	fhirClient, err := fhir.NewFHIRClient(context.Background(), fhir.FHIRConfig{
		ProjectID:       conf.GCLOUD_PROJECT_ID,
		DatasetLocation: conf.GCLOUD_DATASET_LOCATION,
		DatasetID:       conf.GCLOUD_DATASET_ID,
		FHIRStoreID:     conf.GCLOUD_FHIR_STORE_ID,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to setup fhir client:%v", err)
	}
	// external payment service setup
	processor := payment.NewPaymentProcessor(conf.PAYSTACK_API_KEY)
	streamClient, err := streamsdk.NewStreamClient(conf.GETSTREAM_API_KEY, conf.GETSTREAM_API_SECRET)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the getstream client:%v", err)
	}
	// Configutation Options
	opts := server.ConfigOptions{
		Port:                conf.PORT,
		AccessTokenDuration: conf.ACCESS_TOKEN_DURATION,
		AuthMaker:           maker,
		ImageStorage:        imgStorage,
		PaymentProcessor:    processor,
		StreamClient:        streamClient,
		FileStorage:         fileStorage,
		FHIRClient:          fhirClient,
	}
	server := server.NewServer(opts)
	return server, nil
}

func main() {
	server, err := setupServer()
	if err != nil {
		log.Fatalf("fatal error,the server setup process failed : %v", err)
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
