package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	// GET
	r1, err := http.Get("http://www.google.com/robots.txt")
	if err != nil {
		log.Fatalln(err)
	}
	defer r1.Body.Close()
	fmt.Println("GET:", r1.Status)

	// HEAD
	r2, err := http.Head("http://www.google.com/robots.txt")
	if err != nil {
		log.Fatalln(err)
	}
	defer r2.Body.Close()
	fmt.Println("HEAD:", r2.Status)

	// POST
	form := url.Values{}
	form.Add("foo", "bar")
	r3, err := http.Post(
		"http://www.google.com/robots.txt",
		"application/x-www-form-urlencoded",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		log.Fatalln(err)
	}
	defer r3.Body.Close()
	fmt.Println("POST:", r3.Status)

	// PostForm
	form2 := url.Values{}
	form2.Add("foo", "bar")
	r4, err := http.PostForm("http://www.google.com/robots.txt", form2)
	if err != nil {
		log.Fatalln(err)
	}
	defer r4.Body.Close()
	fmt.Println("PostForm:", r4.Status)

	// DELETE
	var client http.Client
	req, err := http.NewRequest("DELETE", "http://www.google.com/robots.txt", nil)
	if err != nil {
		log.Fatalln(err)
	}
	r5, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer r5.Body.Close()
	fmt.Println("DELETE:", r5.Status)

	// PUT
	form3 := url.Values{}
	form3.Add("foo", "bar")
	req2, err := http.NewRequest(
		"PUT",
		"http://www.google.com/robots.txt",
		strings.NewReader(form3.Encode()),
	)
	if err != nil {
		log.Fatalln(err)
	}
	r6, err := client.Do(req2)
	if err != nil {
		log.Fatalln(err)
	}
	defer r6.Body.Close()
	fmt.Println("PUT:", r6.Status)
}
