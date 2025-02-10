package server

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/mbeka02/lyra_backend/internal/auth"
	"github.com/mbeka02/lyra_backend/internal/database"
	"github.com/mbeka02/lyra_backend/internal/objstore"
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
	User *UserHandler
}
type Services struct {
	User service.UserService
}
type Repositories struct {
	User repository.UserRepository
}

func initRepositories(store *database.Store) Repositories {
	return Repositories{
		User: repository.NewUserRepository(store),
	}
}

func initServices(repos Repositories, maker auth.Maker, objStorage objstore.Storage, duration time.Duration) Services {
	return Services{
		User: service.NewUserService(repos.User, maker, objStorage, duration),
	}
}

func initHandlers(services Services) Handlers {
	return Handlers{
		User: NewUserHandler(services.User),
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
