// simple-api-server.go
package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"Message": "All is good with the world",
			"Status":  "Success",
		})
	})
	http.ListenAndServe(":8080", nil)
}
