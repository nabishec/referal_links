package logger

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog/log"
)

func New() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log := log.With().Str("component", "middleware/logger").Logger()

		log.Info().Msg("logger middleware enabled")

		fnc := func(w http.ResponseWriter, r *http.Request) {
			entry := log.With().Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("remote_addr", r.RemoteAddr).
				Str("user_agent", r.UserAgent()).
				Str("request_id", middleware.GetReqID(r.Context())).Logger()

			writer := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			start := time.Now()

			defer func() {
				entry.Info().Int("status", writer.Status()).
					Int("bytes", writer.BytesWritten()).
					Str("duration", time.Since(start).String())
			}()

			next.ServeHTTP(writer, r)
		}
		return http.HandlerFunc(fnc)
	}
}
