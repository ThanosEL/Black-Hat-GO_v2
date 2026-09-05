package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Basic route, restricted to GET requests only.
	r.HandleFunc("/foo", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprint(w, "hi foo")
	}).Methods("GET")

	// Path parameter: matches anything after /users/ and captures it
	// as "user", accessible via mux.Vars(req).
	r.HandleFunc("/users/{user}", func(w http.ResponseWriter, req *http.Request) {
		user := mux.Vars(req)["user"]
		fmt.Fprintf(w, "hi %s\n", user)
	}).Methods("GET")

	// Same idea, but constrained by a regex: only lowercase letters allowed.
	// Anything that doesn't match (e.g. "bob1") returns 404 automatically.
	r.HandleFunc("/strict-users/{user:[a-z]+}", func(w http.ResponseWriter, req *http.Request) {
		user := mux.Vars(req)["user"]
		fmt.Fprintf(w, "hi %s\n", user)
	}).Methods("GET")

	http.ListenAndServe(":8000", r)
}
