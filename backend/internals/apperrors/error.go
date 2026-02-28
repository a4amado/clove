package apperrors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Ensure AppError satisfies the swaggest ErrWithHTTPStatus interface.
var _ interface{ HTTPStatus() int } = (*AppError)(nil)

// AppError represents a structured application error
type AppError struct {
	ID         uuid.UUID `json:"request_id"`
	Type       string    `json:"type"`
	Code       string    `json:"code"`
	Message    string    `json:"message"`
	StatusCode int       `json:"status_code"`
}

func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Code
}

func (e *AppError) HTTPStatus() int {
	return e.StatusCode
}

// WriteError writes an error response to http.ResponseWriter
func WriteError(w http.ResponseWriter, err *AppError) {

	if err.Code == "" {
		err.Code = "INTERNAL_SERVER_ERROR"
	}

	w.WriteHeader(err.StatusCode)
	json.NewEncoder(w).Encode(err)
}

func WriteWsError(c *websocket.Conn, lock *sync.Mutex, err *AppError) {
	lock.Lock()
	defer lock.Unlock()
	c.WriteJSON(err)
	closeMsg := websocket.FormatCloseMessage(websocket.CloseInternalServerErr, err.Code)
	c.WriteMessage(websocket.CloseNormalClosure, closeMsg)
	c.Close()

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

// Sentinel errors — used with usecase.SetExpectedErrors to generate OpenAPI response schemas.
// Each covers a distinct HTTP status code; the AppError schema applies to all.
var (
	ErrUnauthorized        = register("UNAUTHORIZED", "Unauthorized", http.StatusUnauthorized)
	ErrBadRequest          = register("BAD_REQUEST", "Bad request", http.StatusBadRequest)
	ErrNotFound            = register("NOT_FOUND", "Not found", http.StatusNotFound)
	ErrConflict            = register("CONFLICT", "Conflict", http.StatusConflict)
	ErrInternalServerError = register("INTERNAL_SERVER_ERROR", "Internal server error", http.StatusInternalServerError)
)

// GetAllErrorCodes returns all registered error codes (useful for documentation)
func GetAllErrorCodes() []string {
	codes := make([]string, 0, len(errorRegistry))
	for code := range errorRegistry {
		codes = append(codes, code)
	}
	return codes
}
