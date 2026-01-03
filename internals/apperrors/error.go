package apperrors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// AppError represents a structured application error
type AppError struct {
	ID         uuid.UUID     `json:"error_id"`
	Code       string        `json:"code"`
	Message    string        `json:"message"`
	StatusCode int           `json:"status_code"`
	Internal   error         `json:"error"` // The underlying error for logging
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
func WriteError(w *http.ResponseWriter, err *AppError) {
	Log(err)
	if w != nil {
		_ = json.NewEncoder(*w).Encode(err)
	}
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

// GetAllErrorCodes returns all registered error codes (useful for documentation)
func GetAllErrorCodes() []string {
	codes := make([]string, 0, len(errorRegistry))
	for code := range errorRegistry {
		codes = append(codes, code)
	}
	return codes
}
