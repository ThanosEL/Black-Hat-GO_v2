package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// contextKey is a dedicated type for context keys, so our "username" key
// can never accidentally collide with a plain string key used by some
// other package. Using a bare string here (as the book does) triggers a
// go vet warning (SA1029) in modern Go.
type contextKey string

const usernameKey contextKey = "username"

// badAuth is intentionally simplistic — real auth should never compare
// credentials via URL query parameters or plain equality like this. It's
// here purely to demonstrate the negroni.Handler middleware pattern.
type badAuth struct {
	Username string
	Password string
}

func (b *badAuth) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")

	if username != b.Username || password != b.Password {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return // stop the chain here — never call next()
	}

	ctx := context.WithValue(r.Context(), usernameKey, username)
	r = r.WithContext(ctx)
	next(w, r)
}

func hello(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(usernameKey).(string)
	fmt.Fprintf(w, "Hi %s\n", username)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/hello", hello).Methods("GET")

	n := negroni.Classic()
	n.Use(&badAuth{
		Username: "admin",
		Password: "password",
	})
	n.UseHandler(r)

	http.ListenAndServe(":8000", n)
}
