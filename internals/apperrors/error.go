package apperrors

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

// AppError represents a structured application error
type AppError struct {
	ID         uuid.UUID     `json:"id"`
	Code       string        `json:"code"`
	Message    string        `json:"message"`
	StatusCode int           `json:"status_code"`
	Internal   error         `json:"-"` // The underlying error for logging
	Request    *http.Request `json:"-"`
}

func (e *AppError) Error() string {
	if e.Internal != nil {
		return fmt.Sprintf("%s: %s (internal: %v)", e.Code, e.Message, e.Internal)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func Log(err *AppError) {

}

// WriteError writes an error response to http.ResponseWriter
func WriteError(w http.ResponseWriter, err *AppError) {
	Log(err)
	http.Error(w, err.Code, err.StatusCode)
}

// =============================================================================
// ERROR REGISTRY - All application errors defined here
// =============================================================================

// This map ensures no duplicate error codes at compile time
// If you try to register the same code twice, the app will panic on init
var errorRegistry = make(map[string]*AppError)

func register(code, message string, statusCode int) *AppError {
	if _, exists := errorRegistry[code]; exists {
		panic(fmt.Sprintf("duplicate error code registered: %s", code))
	}
	err := &AppError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
	errorRegistry[code] = err
	return err
}

// GetAllErrorCodes returns all registered error codes (useful for documentation)
func GetAllErrorCodes() []string {
	codes := make([]string, 0, len(errorRegistry))
	for code := range errorRegistry {
		codes = append(codes, code)
	}
	return codes
}
