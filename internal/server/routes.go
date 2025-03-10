package server

import (
	"encoding/json"
	"log"
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

	r.Use(httprate.LimitByIP(100, time.Minute))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Use(middleware.Logger)

	r.Get("/", s.TestHandler)
	// TODO: move this to the right place

	r.Get("/health", s.healthHandler)
	r.Post("/register", s.handlers.User.HandleCreateUser)
	r.Post("/login", s.handlers.User.HandleLogin)
	// TODO: Restructure this
	r.Route("/user", func(r chi.Router) {
		r.Use(m.AuthMiddleware(s.opts.AuthMaker))
		r.Get("/", s.handlers.User.HandleGetUser)
		r.Patch("/", s.handlers.User.HandleUpdateUser)
		r.Patch("/profilePicture", s.handlers.User.HandleProfilePicture)

		r.Post("/patient", s.handlers.Patient.HandleCreatePatient)
		r.Post("/appointment", s.handlers.Appointment.HandleCreateAppointment)
		r.Get("/doctor", s.handlers.Doctor.HandleGetDoctors)
		r.Get("/doctor/availability", s.handlers.Availability.HandleGetAvailabilityByDoctor)
		r.Post("/doctor/slots", s.handlers.Availability.HandleGetSlots)

		r.Post("/doctor", s.handlers.Doctor.HandleCreateDoctor)
		r.Post("/doctor/availability", s.handlers.Availability.HandleCreateAvailability)
		r.Delete("/doctor/availability/id/{availabilityId}", s.handlers.Availability.HandleDeleteById)
		r.Delete("/doctor/availability/day/{dayOfWeek}", s.handlers.Availability.HandleDeleteByDay)
	})
	return r
}

func (s *Server) TestHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Lyra API"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}
