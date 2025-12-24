package apierror

type ApiError string
type ApiErrorResponse struct {
	Error   ApiError `json:"error"`
	Message string   `json:"message"`
}

const (
	ErrorCodeInvalidUUID           ApiError = "invalid_id"
	ErrorCodeInvalidRequestBody    ApiError = "invalid_body"
	ErrorCodeUserNotFound          ApiError = "user_not_found"
	ErrorCodeOperationFailed       ApiError = "operation_failed"
	ErrorCodeBodySizeLimitExceeded ApiError = "body_size_limit_exceeded"
)
