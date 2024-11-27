//go:build !solution

package requestlog

import (
	"math/rand"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type CustomResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *CustomResponseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *CustomResponseWriter) Write(b []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

func (rw *CustomResponseWriter) StatusCode() int {
	if rw.statusCode == 0 {
		return http.StatusOK
	}
	return rw.statusCode
}

func keyGeneration() string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func Log(l *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := keyGeneration()
			l.Info("request started",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("request_id", requestID),
			)

			customRW := &CustomResponseWriter{ResponseWriter: w}
			startTime := time.Now()
			defer func() {
				duration := time.Since(startTime)
				if err := recover(); err != nil {
					l.Error("request panicked",
						zap.String("method", r.Method),
						zap.String("path", r.URL.Path),
						zap.String("request_id", requestID),
						zap.Any("error", err),
					)
					panic(err)
				}

				l.Info("request finished",
					zap.String("method", r.Method),
					zap.String("path", r.URL.Path),
					zap.String("request_id", requestID),
					zap.Duration("duration", duration),
					zap.Int("status_code", customRW.StatusCode()),
				)
			}()
			next.ServeHTTP(customRW, r)
		})
	}
}
