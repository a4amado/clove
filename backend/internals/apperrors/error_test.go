package apperrors

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppError_Structure(t *testing.T) {
	err := &AppError{
		ID:         uuid.New(),
		Code:       "TEST_ERROR",
		Message:    "This is a test error",
		StatusCode: http.StatusBadRequest,
	}

	assert.NotEqual(t, uuid.Nil, err.ID)
	assert.Equal(t, "TEST_ERROR", err.Code)
	assert.Equal(t, "This is a test error", err.Message)
	assert.Equal(t, http.StatusBadRequest, err.StatusCode)
}

func TestAppError_JSON(t *testing.T) {
	id := uuid.New()
	err := &AppError{
		ID:         id,
		Code:       "TEST_ERROR",
		Message:    "Test message",
		StatusCode: http.StatusInternalServerError,
	}

	data, jsonErr := json.Marshal(err)
	require.NoError(t, jsonErr)

	var decoded AppError
	jsonErr = json.Unmarshal(data, &decoded)
	require.NoError(t, jsonErr)

	assert.Equal(t, err.ID, decoded.ID)
	assert.Equal(t, err.Code, decoded.Code)
	assert.Equal(t, err.Message, decoded.Message)
	assert.Equal(t, err.StatusCode, decoded.StatusCode)
}

func TestWriteError(t *testing.T) {
	tests := []struct {
		name           string
		err            *AppError
		expectedStatus int
		expectedCode   string
	}{
		{
			name: "writes error with status code",
			err: &AppError{
				ID:         uuid.New(),
				Code:       "BAD_REQUEST",
				Message:    "Invalid input",
				StatusCode: http.StatusBadRequest,
			},
			expectedStatus: http.StatusBadRequest,
			expectedCode:   "BAD_REQUEST",
		},
		{
			name: "writes internal server error",
			err: &AppError{
				ID:         uuid.New(),
				Code:       "INTERNAL_ERROR",
				Message:    "Something went wrong",
				StatusCode: http.StatusInternalServerError,
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   "INTERNAL_ERROR",
		},
		{
			name: "uses default code when empty",
			err: &AppError{
				ID:         uuid.New(),
				Code:       "",
				Message:    "No code provided",
				StatusCode: http.StatusInternalServerError,
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   "INTERNAL_SERVER_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			WriteError(w, tt.err)

			assert.Equal(t, tt.expectedStatus, w.Code)

			// Parse the response body
			var response AppError
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedCode, response.Code)
		})
	}
}

func TestWriteError_ResponseFormat(t *testing.T) {
	w := httptest.NewRecorder()
	id := uuid.New()

	err := &AppError{
		ID:         id,
		Code:       "NOT_FOUND",
		Message:    "Resource not found",
		StatusCode: http.StatusNotFound,
	}

	WriteError(w, err)

	// Verify the response body is valid JSON with expected fields
	var response map[string]interface{}
	jsonErr := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, jsonErr)

	assert.Equal(t, id.String(), response["request_id"])
	assert.Equal(t, "NOT_FOUND", response["code"])
	assert.Equal(t, "Resource not found", response["message"])
	assert.Equal(t, float64(http.StatusNotFound), response["status_code"])
}

func TestGetAllErrorCodes(t *testing.T) {
	codes := GetAllErrorCodes()

	// The registry should contain at least the errors that were registered at init time
	// This is a basic sanity check
	assert.IsType(t, []string{}, codes)
}

func TestAppError_ZeroValue(t *testing.T) {
	var err AppError

	// Zero value should have nil UUID
	assert.Equal(t, uuid.Nil, err.ID)
	assert.Empty(t, err.Code)
	assert.Empty(t, err.Message)
	assert.Equal(t, 0, err.StatusCode)
}

func TestWriteError_MultipleWrites(t *testing.T) {
	// Verify writing multiple errors to different recorders works independently
	w1 := httptest.NewRecorder()
	w2 := httptest.NewRecorder()

	err1 := &AppError{
		ID:         uuid.New(),
		Code:       "ERROR_1",
		Message:    "First error",
		StatusCode: http.StatusBadRequest,
	}

	err2 := &AppError{
		ID:         uuid.New(),
		Code:       "ERROR_2",
		Message:    "Second error",
		StatusCode: http.StatusNotFound,
	}

	WriteError(w1, err1)
	WriteError(w2, err2)

	assert.Equal(t, http.StatusBadRequest, w1.Code)
	assert.Equal(t, http.StatusNotFound, w2.Code)

	var resp1, resp2 AppError
	json.Unmarshal(w1.Body.Bytes(), &resp1)
	json.Unmarshal(w2.Body.Bytes(), &resp2)

	assert.Equal(t, "ERROR_1", resp1.Code)
	assert.Equal(t, "ERROR_2", resp2.Code)
}
