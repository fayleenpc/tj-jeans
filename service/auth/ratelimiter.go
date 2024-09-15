package auth

import (
	"log"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(1, 1)

func WithRateLimiter(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if the request is allowed

		if !limiter.Allow() {
			// Rate limit exceeded
			log.Printf("Rate limit exceeded, %v .\n", http.StatusTooManyRequests)
			// utils.WriteError(w, http.StatusTooManyRequests, fmt.Errorf("Rate limit exceeded"))
			// return
		}

		// Handle the request
		h, m, s := time.Now().Clock()
		log.Printf("Request allowed at %v, %v, %v. \n", h, m, s)
		handlerFunc(w, r)
	}
}
