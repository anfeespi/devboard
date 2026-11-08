// Package middlewares middlewares
package middlewares

import (
	"log/slog"
	"net/http"
	"time"
)

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

// NewResponseRecorder new response recorder
func newResponseRecorder(writer http.ResponseWriter) *responseRecorder {
	return &responseRecorder{
		ResponseWriter: writer,
		statusCode:     http.StatusOK,
	}
}

func (rr *responseRecorder) WriteHeader(code int) {
	if !rr.written {
		rr.statusCode = code
		rr.written = true
		rr.ResponseWriter.WriteHeader(rr.statusCode)
	}
}

func (rr *responseRecorder) Write(b []byte) (int, error) {
	if !rr.written {
		rr.WriteHeader(http.StatusOK)
	}

	return rr.ResponseWriter.Write(b)
}

// Logger logger
func Logger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			start := time.Now()
			rec := newResponseRecorder(writer)

			next.ServeHTTP(rec.ResponseWriter, request)

			logger.InfoContext(request.Context(), "request completed",
				slog.String("method", request.Method),
				slog.String("path", request.URL.Path),
				slog.Int("status", rec.statusCode),
				slog.Duration("duration", time.Since(start)),
				slog.String("remote_addr", request.RemoteAddr),
			)
		})
	}
}
