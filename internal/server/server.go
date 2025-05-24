package server

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/mbeka02/lyra_backend/internal/auth"
	"github.com/mbeka02/lyra_backend/internal/database"
	"github.com/mbeka02/lyra_backend/internal/fhir"
	"github.com/mbeka02/lyra_backend/internal/objstore"
	"github.com/mbeka02/lyra_backend/internal/payment"
	"github.com/mbeka02/lyra_backend/internal/server/handler"
	"github.com/mbeka02/lyra_backend/internal/server/repository"
	"github.com/mbeka02/lyra_backend/internal/server/service"
	"github.com/mbeka02/lyra_backend/internal/streamsdk"
)

type ConfigOptions struct {
	Port                string
	AccessTokenDuration time.Duration
	AuthMaker           auth.Maker
	ImageStorage        objstore.Storage
	FileStorage         objstore.Storage
	PaymentProcessor    *payment.PaymentProcessor
	StreamClient        *streamsdk.StreamClient
	FHIRClient          *fhir.FHIRClient
}
type Server struct {
	opts     ConfigOptions
	db       *database.Store
	handlers Handlers
}
type Handlers struct {
	User              *handler.UserHandler
	Patient           *handler.PatientHandler
	Doctor            *handler.DoctorHandler
	Availability      *handler.AvailabilityHandler
	Appointment       *handler.AppointmentHandler
	Payment           *handler.PaymentHandler
	DocumentReference *handler.DocumentReferenceHandler
}
type Services struct {
	User              service.UserService
	Patient           service.PatientService
	Doctor            service.DoctorService
	Availability      service.AvailabilityService
	Appointment       service.AppointmentService
	Payment           service.PaymentService
	DocumentReference service.DocumentReferenceService
	Observation       service.ObservationService
}
type Repositories struct {
	User         repository.UserRepository
	Patient      repository.PatientRepository
	Doctor       repository.DoctorRepository
	Availability repository.AvailabilityRepository
	Appointment  repository.AppointmentRepository
	Payment      repository.PaymentRepository
}

func initRepositories(store *database.Store) Repositories {
	return Repositories{
		User:         repository.NewUserRepository(store),
		Patient:      repository.NewPatientRepository(store),
		Doctor:       repository.NewDoctorRepository(store),
		Availability: repository.NewAvailabilityRepository(store),
		Appointment:  repository.NewAppointmentRepository(store),
		Payment:      repository.NewPaymentRepository(store),
	}
}

func initServices(repos Repositories, maker auth.Maker, imgStorage, fileStorage objstore.Storage, duration time.Duration, paymentProcessor *payment.PaymentProcessor, streamClient *streamsdk.StreamClient, fhirClient *fhir.FHIRClient) Services {
	return Services{
		User:              service.NewUserService(repos.User, maker, streamClient, imgStorage, duration),
		Patient:           service.NewPatientService(repos.Patient, fhirClient, fileStorage),
		Doctor:            service.NewDoctorService(repos.Doctor, repos.Appointment),
		Availability:      service.NewAvailabilityService(repos.Availability, repos.Doctor),
		Appointment:       service.NewAppointmentService(repos.Appointment, repos.Patient, repos.Doctor, paymentProcessor),
		Payment:           service.NewPaymentService(paymentProcessor, repos.Payment),
		DocumentReference: service.NewDocumentReferenceService(fhirClient, fileStorage),
		Observation:       service.NewObservationService(fhirClient),
	}
}

func initHandlers(services Services) Handlers {
	return Handlers{
		User:              handler.NewUserHandler(services.User),
		Patient:           handler.NewPatientHandler(services.Patient),
		Doctor:            handler.NewDoctorHandler(services.Doctor),
		Availability:      handler.NewAvailabilityHandler(services.Availability),
		Appointment:       handler.NewAppointmentHandler(services.Appointment),
		Payment:           handler.NewPaymentHandler(services.Payment),
		DocumentReference: handler.NewDocumentReferenceHandler(services.Patient, services.Doctor, services.DocumentReference),
	}
}

func NewServer(opts ConfigOptions) *http.Server {
	store := database.NewStore()
	// repository(data access) layer
	repositories := initRepositories(store)
	// service layer
	services := initServices(repositories, opts.AuthMaker, opts.ImageStorage, opts.FileStorage, opts.AccessTokenDuration, opts.PaymentProcessor, opts.StreamClient, opts.FHIRClient)
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
		ReadTimeout:  45 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
