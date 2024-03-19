package main

import (
	"log/slog"
	"net/http"
)

type Middlewares struct {
	logger *slog.Logger
}

func (m Middlewares) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			method = r.Method
			uri    = r.URL.RequestURI()
		)
		m.logger.Info("request received", "ip", ip, "method", method, "uri", uri)

		next.ServeHTTP(w, r)
	})
}

func (m Middlewares) recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				m.logger.Error("the application recovered from panic", "error", err)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (m Middlewares) secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")

		next.ServeHTTP(w, r)
	})
}
