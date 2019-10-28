package middleware

import (
	"log"
	"net/http"
	"time"
)

//Middleware struct
type Middleware struct {
}

// LoggingHanlder  log is time of http requests
func (m Middleware) LoggingHanlder(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		h.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("[%s]  %q %v", r.Method, r.URL.String(), t2.Sub(t1))
	}
	return http.HandlerFunc(fn)
}
