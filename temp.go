package main

import (
	"context"
	"fmt"
	"net/http"
)

// Custom key type to avoid collisions
type key string

const userKey key = "user"

// User struct to simulate a user object
type User struct {
	ID    int
	Email string
	Name  string
}

// Middleware to add a user to the request context
func UserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a user being authenticated
		user := &User{ID: 1, Email: "user@example.com", Name: "John Doe"}

		// Attach user to the request's context
		ctx := context.WithValue(r.Context(), userKey, user)

		// Pass the request with the new context to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Handler that retrieves the user from the context
func HandleRequest(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user from the contextgo 
	user, ok := r.Context().Value(userKey).(*User)
	if !ok {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Respond with the user's information
	fmt.Fprintf(w, "User: %s (Email: %s)\n", user.Name, user.Email)
}

func main() {
	// Setup the routes and middleware
	http.Handle("/", UserMiddleware(http.HandlerFunc(HandleRequest)))

	// Start the server
	http.ListenAndServe(":8080", nil)
}
