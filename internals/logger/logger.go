package logger

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/google/uuid"
)

var (
	sentryHandler *sentryhttp.Handler
)

const (
	PHASE = "phase"
	Error = "error"
)

func init() {

	// Initialize Sentry
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              "https://bca638ba753055605880335d5256d441@o4510686044880896.ingest.de.sentry.io/4510686197448784",
		Environment:      "production", // or "development"
		TracesSampleRate: 1.0,          // Adjust based on traffic (0.0 to 1.0)
		EnableTracing:    true,
	}); err != nil {
		log.Fatalf("Sentry initialization failed: %v\n", err)
	}

	// Create Sentry HTTP handler with options
	sentryHandler = sentryhttp.New(sentryhttp.Options{
		Repanic:         true,  // Re-panic after capturing
		WaitForDelivery: false, // Don't block on event delivery
	})
}

// FlushSentry should be called before application shutdown
func FlushSentry() {
	sentry.Flush(2 * time.Second)
}

// GetSentryHandler returns the initialized Sentry HTTP handler
func GetSentryHandler() *sentryhttp.Handler {
	return sentryHandler
}

type captureResponseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

func (crw *captureResponseWriter) WriteHeader(statusCode int) {
	if !crw.written {
		crw.statusCode = statusCode
		crw.written = true
		crw.ResponseWriter.WriteHeader(statusCode)
	}
}

func (crw *captureResponseWriter) Write(b []byte) (int, error) {
	if !crw.written {
		crw.WriteHeader(http.StatusOK)
	}
	return crw.ResponseWriter.Write(b)
}

// CaptureHTTPErrors middleware captures HTTP errors and sends them to Sentry
func CaptureHTTPErrors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a custom response writer to capture the status code
		crw := &captureResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			written:        false,
		}

		next(crw, r)

		// Capture 4xx and 5xx errors
		if crw.statusCode >= 400 {
			sentry.WithScope(func(scope *sentry.Scope) {
				scope.SetTag("http.status_code", fmt.Sprintf("%d", crw.statusCode))
				scope.SetTag("http.method", r.Method)
				scope.SetTag("http.path", r.URL.Path)

				scope.SetContext("request", map[string]interface{}{
					"url":         r.URL.String(),
					"method":      r.Method,
					"remote_addr": r.RemoteAddr,
					"user_agent":  r.UserAgent(),
				})

				level := sentry.LevelWarning
				if crw.statusCode >= 500 {
					level = sentry.LevelError
				}
				scope.SetLevel(level)

				sentry.CaptureMessage(fmt.Sprintf("HTTP %d on %s %s", crw.statusCode, r.Method, r.URL.Path))
			})
		}
	}
}

// WrapWithSentry wraps a handler with both Sentry and error capturing middleware
func WrapWithSentry(handler http.HandlerFunc) http.HandlerFunc {
	return sentryHandler.HandleFunc(CaptureHTTPErrors(handler))
}

// JSONLoggerParams contains parameters for creating a request logger
type JSONLoggerParams struct {
	RequestID uuid.UUID
	Req       *http.Request
}
