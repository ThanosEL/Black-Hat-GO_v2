package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type logger struct {
	Inner http.Handler
}

func (l *logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("start %s\n", time.Now().String())
	l.Inner.ServeHTTP(w, r)
	log.Println("finish %s")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func main() {
	f := http.HandlerFunc(hello)
	l := logger{Inner: f}
	http.ListenAndServe(":8000", &l)
}
