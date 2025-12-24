package apierror

import (
	"bytes"
	"encoding/json"
)

type ApiError string
type ApiErrorResponse struct {
	Error   ApiError `json:"error"`
	Message string   `json:"message"`
}

func FormatErrorJSON(apiError ApiError, message string) string {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(ApiErrorResponse{
		Error:   apiError,
		Message: message,
	})
	if err != nil {
		panic(err)
	}
	return buf.String()
}

const (
	ErrorCodeInvalidUUID           ApiError = "invalid_id"
	ErrorCodeInvalidRequestBody    ApiError = "invalid_body"
	ErrorCodeUserNotFound          ApiError = "user_not_found"
	ErrorCodeOperationFailed       ApiError = "operation_failed"
	ErrorCodeBodySizeLimitExceeded ApiError = "body_size_limit_exceeded"
)
