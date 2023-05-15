// Http middleware
// pkg/http/middleware/middleware.go
package middleware

import (
	"log"
	"net/http"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Printf("Received request: %s %s", r.Method, r.RequestURI)

		// Call the next handler, which can be another middleware in the chain, or the final handler
		next.ServeHTTP(w, r)

		// Do more stuff here, like logging the response code and time
		log.Printf("Completed request: %s %s", r.Method, r.RequestURI)
	})
}
