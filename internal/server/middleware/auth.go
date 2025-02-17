package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/mbeka02/lyra_backend/internal/auth"
)

type contextKey string

// APIError represents a structured error response
type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

const (
	authorizationTypeBearer            = "Bearer"
	authorizationPayloadKey contextKey = "authorization_payload"
)

var (
	ErrMissingAuth     = errors.New("authorization header is missing")
	ErrMalformedAuth   = errors.New("malformed authorization header")
	ErrUnsupportedAuth = errors.New("unsupported authorization type")
	ErrInvalidPayload  = errors.New("invalid authorization payload")
)

func GetAuthPayload(ctx context.Context) (*auth.Payload, error) {
	payload, ok := ctx.Value(authorizationPayloadKey).(*auth.Payload)
	if !ok {
		return nil, ErrInvalidPayload
	}
	return payload, nil
}

func respondWithVerificationError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	apiError := APIError{
		Status:  http.StatusUnauthorized,
		Message: http.StatusText(http.StatusUnauthorized),
		Detail:  err.Error(),
	}
	json.NewEncoder(w).Encode(apiError)
}

func extractAndVerifyToken(r *http.Request, maker auth.Maker) (*auth.Payload, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, ErrMissingAuth
	}

	fields := strings.Fields(authHeader)
	if len(fields) != 2 {
		return nil, ErrMalformedAuth
	}

	authType := fields[0]
	if authType != authorizationTypeBearer {
		return nil, ErrUnsupportedAuth
	}

	return maker.Verify(fields[1])
}

func AuthMiddleware(maker auth.Maker) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			payload, err := extractAndVerifyToken(r, maker)
			if err != nil {
				respondWithVerificationError(w, err)
				return
			}

			// Create new context with the payload
			ctx := context.WithValue(r.Context(), authorizationPayloadKey, payload)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
