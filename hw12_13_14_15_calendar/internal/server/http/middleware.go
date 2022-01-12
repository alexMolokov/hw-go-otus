package internalhttp

import (
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (s Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := newResponseWriter(w)
		next.ServeHTTP(rw, r)

		ip := r.RemoteAddr
		if v := r.Header.Get("X-FORWARDED-FOR"); v != "" {
			ip = v
		}

		s.Logger.Info("%s IP=%s Method=%s Path=%s Code=%d Latency=%d User-Agent:%s",
			r.Proto, ip, r.Method, r.URL.Path, rw.statusCode, time.Since(start), r.UserAgent(),
		)
	})
}
