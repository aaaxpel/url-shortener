package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/aaaxpel/url-shortener/database"
)

var db = database.Connect()

type Url struct {
	ID          int
	ShortCode   string
	URL         string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	AccessCount int
}

func Create(w http.ResponseWriter, r *http.Request) {

	length := 8
	byteLength := (length + 1) / 2
	randomBytes := make([]byte, byteLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Fatal("boom")
	}

	shortCode := hex.EncodeToString(randomBytes)

	query := "INSERT INTO urls (short_code, url) VALUES (?, ?) RETURNING id, short_code, url, createdAt, updatedAt"
	newURL := Url{}
	err = db.QueryRow(query, shortCode, r.PostFormValue("url")).Scan(&newURL.ID, &newURL.ShortCode, &newURL.URL, &newURL.CreatedAt, &newURL.UpdatedAt)
	if err != nil {
		log.Printf("error: %s", err)
	}

	w.WriteHeader(201)

	new, err := json.Marshal(newURL)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	w.Write([]byte(new))
}

func Get(w http.ResponseWriter, r *http.Request) {
	//
	http.Redirect(w, r, r.URL.Path, 200)
}

func Update(w http.ResponseWriter, r *http.Request) {

}

func Delete(w http.ResponseWriter, r *http.Request) {

}

func GetStats(w http.ResponseWriter, r *http.Request) {

}
