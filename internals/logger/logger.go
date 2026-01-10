package logger

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/google/uuid"
)

var reqLogger *slog.Logger

const (
	PHASE = "phase"
	Error = "error"
)

func init() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	})

	reqLogger = slog.New(handler)
}

// req_id := uuid.New()
// log := logger.ReqLogger().With(
// 	slog.String("req_id", req_id.String()),
// 	slog.String("path", r.URL.RawPath),
// 	slog.String("method", r.Method),
// )

type JSONLoggerParams struct {
	RequestID uuid.UUID
	Req       *http.Request
}

func NewJSONLogger(args JSONLoggerParams) *slog.Logger {
	return reqLogger.With(
		slog.String("req_id", args.RequestID.String()),
		slog.String("method", args.Req.Method),
		slog.String("path", args.Req.URL.String()),
	)
}
