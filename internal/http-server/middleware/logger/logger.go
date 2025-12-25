package logger

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log := log.With(slog.String("component", "middleware/logger"))

		log.Info("logger middleware enabled")

		fn := func(writer http.ResponseWriter, req *http.Request) {
			entry := log.With(
				slog.String("method", req.Method),
				slog.String("path", req.URL.Path),
				slog.String("remote_addr", req.RemoteAddr),
				slog.String("user_agent", req.UserAgent()),
				slog.String("request_id", middleware.GetReqID(req.Context())),
			)
			wrappedResponseWriter := middleware.NewWrapResponseWriter(writer, req.ProtoMajor)

			startTime := time.Now()

			defer func() {
				entry.Info("request completed",
					slog.Int("status", wrappedResponseWriter.Status()),
					slog.Int("bytes", wrappedResponseWriter.BytesWritten()),
					slog.String("duration", time.Since(startTime).String()),
				)
			}()

			next.ServeHTTP(wrappedResponseWriter, req)
		}

		return http.HandlerFunc(fn)
	}
}