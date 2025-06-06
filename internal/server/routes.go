package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	m "github.com/mbeka02/lyra_backend/internal/server/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	// Apply global middleware
	r.Use(httprate.LimitByIP(100, time.Minute))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	// logs requests
	r.Use(middleware.Logger)
	// catches panics in the handlers and returns a 500 instead of crashing the server
	r.Use(middleware.Recoverer)
	// extracts the real client IP from the headers even when behind a proxy
	r.Use(middleware.RealIP)
	// add a unique request ID for each request
	r.Use(middleware.RequestID)
	// request timeout
	r.Use(middleware.Timeout(30 * time.Second))
	// Public routes that don't require authentication
	r.Get("/", s.TestHandler)
	r.Get("/health", s.healthHandler)
	r.Post("/register", s.handlers.User.HandleCreateUser)
	r.Post("/login", s.handlers.User.HandleLogin)

	// API versioning - all API endpoints under /api/v1
	r.Route("/api/v1", func(r chi.Router) {
		// Payment endpoints
		r.Route("/payments", func(r chi.Router) {
			r.Post("/webhook", s.handlers.Payment.PaymentWebhook)
			r.Get("/callback", s.handlers.Payment.PaymentCallback)
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			// Apply authentication middleware to all routes in this group
			r.Use(m.AuthMiddleware(s.opts.AuthMaker))

			// User endpoints
			r.Route("/users", func(r chi.Router) {
				r.Get("/me", s.handlers.User.HandleGetUser)
				r.Patch("/me", s.handlers.User.HandleUpdateUser)
				r.Patch("/me/profile-picture", s.handlers.User.HandleProfilePicture)
			})

			// Patient endpoints
			r.Route("/patients", func(r chi.Router) {
				r.Post("/", s.handlers.Patient.HandleCreatePatient)
				r.Get("/appointments", s.handlers.Appointment.HandleGetPatientAppointments)
				r.Get("/{patientId}", s.handlers.Patient.HandleGetPatient)
				// allergies
				r.Route("/{patientId}/allergies", func(r chi.Router) {
					r.Post("/", s.handlers.Allergy.HandleCreateAllergy)
					r.Get("/", s.handlers.Allergy.HandleListAllergies)
					r.Route("/{allergyId}", func(r chi.Router) {
						r.Get("/", s.handlers.Allergy.HandleGetAllergy)
						r.Put("/", s.handlers.Allergy.HandleUpdateAllergy)
						r.Delete("/", s.handlers.Allergy.HandleDeleteAllergy)
					})
				})
				// medications
				r.Route("/{patientId}/medications", func(r chi.Router) {
					r.Post("/", s.handlers.MedicationStatement.HandleCreateMedication)
					r.Get("/", s.handlers.MedicationStatement.HandleListMedications)
					r.Route("/{medicationId}", func(r chi.Router) {
						r.Get("/", s.handlers.MedicationStatement.HandleGetMedication)
						r.Put("/", s.handlers.MedicationStatement.HandleUpdateMedication)
						r.Delete("/", s.handlers.MedicationStatement.HandleDeleteMedication)
					})
				})
			})
			// Observations
			r.Route("/{patientId}/observations", func(r chi.Router) {
				r.Post("/", s.handlers.Observation.HandleCreateObservationInDB)
				r.Get("/", s.handlers.Observation.HandleListObservations)
				r.Route("/{observationId}", func(r chi.Router) {
					r.Get("/", s.handlers.Observation.HandleGetObservation)
					r.Put("/", s.handlers.Observation.HandleUpdateObservation)
					r.Delete("/", s.handlers.Observation.HandleDeleteObservation)
				})
			})
			// Doctor endpoints
			r.Route("/doctors", func(r chi.Router) {
				r.Get("/", s.handlers.Doctor.HandleGetDoctors)
				r.Get("/my-patients", s.handlers.Doctor.HandleListMyPatients)
				r.Post("/", s.handlers.Doctor.HandleCreateDoctor)
				r.Get("/appointments", s.handlers.Appointment.HandleGetDoctorAppointments)

				// Doctor availability endpoints
				r.Route("/availability", func(r chi.Router) {
					r.Get("/", s.handlers.Availability.HandleGetAvailabilityByDoctor)
					r.Post("/", s.handlers.Availability.HandleCreateAvailability)
					r.Post("/slots", s.handlers.Availability.HandleGetSlots)
					r.Delete("/id/{availabilityId}", s.handlers.Availability.HandleDeleteById)
					r.Delete("/day/{dayOfWeek}", s.handlers.Availability.HandleDeleteByDay)
				})
			})

			// Appointment endpoints
			r.Route("/appointments", func(r chi.Router) {
				r.Patch("/status", s.handlers.Appointment.HandleUpdateStatus)
				r.Get("/completed", s.handlers.Appointment.HandleGetCompletedAppointments)
				r.Post("/", s.handlers.Appointment.HandleCreateAppointment)
			})
			// Document endpoints
			r.Route("/documents", func(r chi.Router) {
				r.Get("/", s.handlers.DocumentReference.HandleListPatientDocuments)
				r.Post("/upload", s.handlers.DocumentReference.HandleCreateDocumentReference)
				// Observation endpoints
				r.Route("/observations", func(r chi.Router) {
					r.Get("/", s.handlers.Observation.HandleListPatientObservations)
					r.Post("/note", s.handlers.Observation.HandleCreateConsultationNote)
				})
				// Endpoint to create a signed URL for document operations
				r.Post("/signed-url", s.handlers.DocumentReference.HandleCreateSignedURL)
			})
		})
	})

	return r
}

func (s *Server) TestHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Lyra API"
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, err := json.Marshal(s.db.Health())
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Write(jsonResp)
}
