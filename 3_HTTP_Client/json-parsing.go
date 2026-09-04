package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Status struct {
	Message string
	Status  string
}

func main() {
	res, err := http.Post(
		"http://localhost:8080/ping",
		"application/json",
		nil,
	)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	var status Status
	if err := json.NewDecoder(res.Body).Decode(&status); err != nil {
		log.Fatalln(err)
	}

	log.Printf("%s -> %s\n", status.Status, status.Message)
}
