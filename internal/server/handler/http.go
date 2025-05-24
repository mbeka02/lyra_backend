package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/mbeka02/lyra_backend/internal/auth"
	"github.com/mbeka02/lyra_backend/internal/server/middleware"
)

var validate *validator.Validate

// APIError represents a structured error response
type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

// APIResponse represents a successful response
type APIResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ValidationError represents validation error details
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// QueryParamExtractor provides a generic way to extract and parse query parameters
type QueryParamExtractor struct {
	query url.Values
}

// NewQueryParamExtractor creates a new extractor from an HTTP request
func NewQueryParamExtractor(r *http.Request) *QueryParamExtractor {
	return &QueryParamExtractor{
		query: r.URL.Query(),
	}
}

// Extract a  string parameter with an optional default value
func (q *QueryParamExtractor) GetString(key string, defaultVal ...string) string {
	value := q.query.Get(key)
	if value == "" && len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return value
}

// GetInt extracts an integer parameter with an optional default value
func (q *QueryParamExtractor) GetInt(key string, defaultVal ...int) int {
	value := q.query.Get(key)
	if value == "" && len(defaultVal) > 0 {
		return defaultVal[0]
	}

	result, err := strconv.Atoi(value)
	if err != nil && len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return result
}

// GetInt32 extracts an int32 parameter with an optional default value
func (q *QueryParamExtractor) GetInt32(key string, defaultVal ...int32) int32 {
	value := q.query.Get(key)
	if value == "" && len(defaultVal) > 0 {
		return defaultVal[0]
	}

	result, err := strconv.ParseInt(value, 10, 32)
	if err != nil && len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return int32(result)
}

// GetInt64 extracts an int64 parameter with an optional default value
func (q *QueryParamExtractor) GetInt64(key string, defaultVal ...int64) int64 {
	value := q.query.Get(key)
	if value == "" && len(defaultVal) > 0 {
		return defaultVal[0]
	}

	result, err := strconv.ParseInt(value, 10, 64)
	if err != nil && len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return int64(result)
}

// GetFloat64 extracts a float64 parameter with an optional default value
func (q *QueryParamExtractor) GetFloat64(key string, defaultVal ...float64) float64 {
	value := q.query.Get(key)
	if value == "" && len(defaultVal) > 0 {
		return defaultVal[0]
	}

	result, err := strconv.ParseFloat(value, 64)
	if err != nil && len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return result
}

// GetBool extracts a boolean parameter with an optional default value
func (q *QueryParamExtractor) GetBool(key string, defaultVal ...bool) bool {
	value := q.query.Get(key)
	if value == "" && len(defaultVal) > 0 {
		return defaultVal[0]
	}

	result, err := strconv.ParseBool(value)
	if err != nil && len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return result
}

var ErrInvalidJSON = errors.New("invalid JSON payload")

// respondWithJSON handles writing JSON responses
func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		return fmt.Errorf("failed to encode response: %w", err)
	}
	return nil
}

func respondWithImage(w http.ResponseWriter, data []byte) error {
	w.Header().Set("Content-Type", "image/jpeg")
	if _, err := w.Write(data); err != nil {
		return fmt.Errorf("unable to write the data to the connection:%v", err)
	}
	return nil
}

// respondWithError handles error responses in a consistent format
func respondWithError(w http.ResponseWriter, status int, err error) {
	apiError := APIError{
		Status:  status,
		Message: http.StatusText(status),
		Detail:  err.Error(),
	}

	respondWithJSON(w, status, apiError)
}

// parseJSON safely decodes JSON request bodies
func parseJSON(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidJSON, err)
	}

	return nil
}

// validateRequest handles struct validation
func validateRequest(v interface{}) []ValidationError {
	if err := validate.Struct(v); err != nil {
		var validationErrors []ValidationError
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, ValidationError{
				Field:   err.Field(),
				Message: fmt.Sprintf("failed validation on '%s'", err.Tag()),
			})
		}
		return validationErrors
	}
	return nil
}

func parseAndValidateRequest(r *http.Request, v interface{}) error {
	if err := parseJSON(r, v); err != nil {
		return err
	}
	if validationErrors := validateRequest(v); validationErrors != nil {
		return fmt.Errorf("Validation failed : %v", validationErrors)
	}
	return nil
}

// getAuthPayload extracts the auth payload from the request context.
// If extraction fails, it writes an error response and returns false.
// this is a wrapper function for the GetAuthPayload() method for the middleware package
func getAuthPayload(w http.ResponseWriter, r *http.Request) (*auth.Payload, bool) {
	payload, err := middleware.GetAuthPayload(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return nil, false
	}
	return payload, true
}

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}
