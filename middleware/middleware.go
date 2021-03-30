package middleware

import (
	"log"
	"net/http"
	"time"
)

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func Timer(handler http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		recorder := &StatusRecorder{
			ResponseWriter: w,
			Status:         200,
		}
		handler.ServeHTTP(recorder, r)
		duration := time.Now().Sub(startTime)

		log.Printf("[Log] %s | %3d | %13v | %s\n",
			r.Method,
			recorder.Status,
			duration.String(),
			r.URL.Path,
		)
	})
}
