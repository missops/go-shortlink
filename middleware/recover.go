package middleware

import (
	"log"
	"net/http"
)

//RecoverHanlder recover for panic 
func (m Middleware) RecoverHanlder(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request)  {
		defer func ()  {
			if err := recover(); err != nil {
				log.Printf("recover from panic %v", err)
				http.Error(w, http.StatusText(500),500)
			}
		}()
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}