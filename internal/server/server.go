package server

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/mbeka02/lyra_backend/internal/auth"
	"github.com/mbeka02/lyra_backend/internal/database"
	"github.com/mbeka02/lyra_backend/internal/objstore"
	"github.com/mbeka02/lyra_backend/internal/server/handler"
	"github.com/mbeka02/lyra_backend/internal/server/repository"
	"github.com/mbeka02/lyra_backend/internal/server/service"
)

type ConfigOptions struct {
	Port                string
	AccessTokenDuration time.Duration
	AuthMaker           auth.Maker
	ObjectStorage       objstore.Storage
}
type Server struct {
	opts     ConfigOptions
	db       *database.Store
	handlers Handlers
}
type Handlers struct {
	User         *handler.UserHandler
	Patient      *handler.PatientHandler
	Doctor       *handler.DoctorHandler
	Availability *handler.AvailabilityHandler
}
type Services struct {
	User         service.UserService
	Patient      service.PatientService
	Doctor       service.DoctorService
	Availability service.AvailabilityService
}
type Repositories struct {
	User         repository.UserRepository
	Patient      repository.PatientRepository
	Doctor       repository.DoctorRepository
	Availability repository.AvailabilityRepository
}

func initRepositories(store *database.Store) Repositories {
	return Repositories{
		User:         repository.NewUserRepository(store),
		Patient:      repository.NewPatientRepository(store),
		Doctor:       repository.NewDoctorRepository(store),
		Availability: repository.NewAvailabilityRepository(store),
	}
}

func initServices(repos Repositories, maker auth.Maker, objStorage objstore.Storage, duration time.Duration) Services {
	return Services{
		User:         service.NewUserService(repos.User, maker, objStorage, duration),
		Patient:      service.NewPatientService(repos.Patient),
		Doctor:       service.NewDoctorService(repos.Doctor),
		Availability: service.NewAvailabilityService(repos.Availability, repos.Doctor),
	}
}

func initHandlers(services Services) Handlers {
	return Handlers{
		User:         handler.NewUserHandler(services.User),
		Patient:      handler.NewPatientHandler(services.Patient),
		Doctor:       handler.NewDoctorHandler(services.Doctor),
		Availability: handler.NewAvailabilityHandler(services.Availability),
	}
}

func NewServer(opts ConfigOptions) *http.Server {
	store := database.NewStore()
	// repository(data access) layer
	repositories := initRepositories(store)
	// service layer
	services := initServices(repositories, opts.AuthMaker, opts.ObjectStorage, opts.AccessTokenDuration)
	// transport layer
	handlers := initHandlers(services)

	NewServer := &Server{
		opts:     opts,
		db:       store,
		handlers: handlers,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", NewServer.opts.Port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
