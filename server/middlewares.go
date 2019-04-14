package server

import (
	"log"
	"net/http"
	"time"
)

func withLogs(next http.Handler) http.Handler {
	return http.HandlerFunc(func(wr http.ResponseWriter, req *http.Request) {
		start := time.Now()
		defer func() {
			dur := time.Now().Sub(start)
			log.Printf("completed {Method=%s, Path=%s} in %s", req.Method, req.URL.Path, dur.String())
		}()

		next.ServeHTTP(wr, req)
	})
}
