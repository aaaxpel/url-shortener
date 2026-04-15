package main

import (
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("POST /shorten", Create)
	http.HandleFunc("GET /shorten/{code}", Get)
	http.HandleFunc("PUT /shorten/{code}", Update)
	http.HandleFunc("DELETE /shorten/{code}", Delete)
	http.HandleFunc("GET /shorten/{code}/stats", GetStats)

	// http.Redirect()

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
