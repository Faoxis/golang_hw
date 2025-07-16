package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

func loggingMiddleware(logger Logger, next http.Handler) http.Handler { //nolint:unused
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Создаем обертку над ResponseWriter, чтобы перехватить статус
		ww := &responseWriter{ResponseWriter: w, statusCode: 200}

		next.ServeHTTP(ww, r)

		duration := time.Since(start)

		clientIP := r.RemoteAddr
		if ip := r.Header.Get("X-Real-IP"); ip != "" {
			clientIP = ip
		}

		userAgent := r.UserAgent()
		method := r.Method
		path := r.URL.RequestURI()
		proto := r.Proto
		status := ww.statusCode
		timestamp := time.Now().Format("02/Jan/2006:15:04:05 -0700")

		logLine := fmt.Sprintf(`%s [%s] %s %s %s %d %v "%s"`,
			clientIP, timestamp, method, path, proto, status, duration.Milliseconds(), userAgent)

		logger.Info(logLine)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}
