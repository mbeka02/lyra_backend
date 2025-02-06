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
)

type Server struct {
	port int

	db          *database.Store
	UserHandler *UserHandler
}

func NewServer(maker auth.Maker, duration time.Duration) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	store := database.NewStore()
	NewServer := &Server{
		port:        port,
		db:          store,
		UserHandler: &UserHandler{Store: store, AuthMaker: maker, AccessTokenDuration: duration},
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
