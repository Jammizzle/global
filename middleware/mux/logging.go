package mux

import (
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type statusWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		path := r.RequestURI
		sw := statusWriter{ResponseWriter: w}

		next.ServeHTTP(w, r)

		end := time.Now()
		latency := end.Sub(start)

		entry := logrus.WithFields(logrus.Fields{
			"status":  sw.status,
			"method":  r.Method,
			"path":    path,
			"ip":      r.RemoteAddr,
			"latency": latency,
			"time":    end.Format(time.RFC3339),
		})

		if os.Getenv("DEBUG") == "true" {
			entry.Info("HTTP Request")
		}
	})
}
