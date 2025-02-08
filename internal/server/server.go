package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/mbeka02/lyra_backend/internal/auth"
	"github.com/mbeka02/lyra_backend/internal/database"
	"github.com/mbeka02/lyra_backend/internal/imgstore"
	"github.com/mbeka02/lyra_backend/internal/server/repository"
	"github.com/mbeka02/lyra_backend/internal/server/service"
)

type Server struct {
	port     int
	db       *database.Store
	handlers Handlers
}
type Handlers struct {
	User *UserHandler
}
type Services struct {
	User service.UserService
}
type Repositories struct {
	User repository.UserRepository
}

func initRepositores(store *database.Store) Repositories {
	return Repositories{
		User: repository.NewUserRepository(store),
	}
}

func initServices(repos Repositories, maker auth.Maker, imgStorage imgstore.Storage, duration time.Duration) Services {
	return Services{
		User: service.NewUserService(repos.User, maker, imgStorage, duration),
	}
}

func initHandlers(services Services) Handlers {
	return Handlers{
		User: NewUserHandler(services.User),
	}
}

func NewServer(maker auth.Maker, imgStorage imgstore.Storage, duration time.Duration) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	store := database.NewStore()
	// repository(data access) layer
	repositories := initRepositores(store)
	// service layer
	services := initServices(repositories, maker, imgStorage, duration)
	// transport layer
	handlers := initHandlers(services)
	NewServer := &Server{
		port:     port,
		db:       store,
		handlers: handlers,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
