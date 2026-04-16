package main

import (
	"crypto/rand"
	"database/sql"
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
	// Doing it in 1 row doesn't guarantee that it's bad request or bad database query / server issue
	// Should ideally split it up to account for that
	err = db.QueryRow(query, shortCode, r.PostFormValue("url")).Scan(&newURL.ID, &newURL.ShortCode, &newURL.URL, &newURL.CreatedAt, &newURL.UpdatedAt)
	if err != nil {
		log.Printf("error: %s", err)
		error, _ := json.Marshal(err)
		w.WriteHeader(400)
		w.Write([]byte(error))
	}

	urlJSON, err := json.Marshal(newURL)
	if err != nil {
		log.Printf("Error: %s", err)
		error, _ := json.Marshal(err)
		w.WriteHeader(500)
		w.Write([]byte(error))
		return
	}

	w.WriteHeader(201)
	w.Write([]byte(urlJSON))
}

func Get(w http.ResponseWriter, r *http.Request) {

	query := "SELECT * FROM urls WHERE short_code = ?"
	url := Url{}
	err := db.QueryRow(query, r.PathValue("code")).Scan(&url.ID, &url.ShortCode, &url.URL, &url.CreatedAt, &url.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return

			// Alternative
			// http.NotFound(w, r)
		}
	}

	urlJSON, err := json.Marshal(url)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	w.WriteHeader(200)
	w.Write([]byte(urlJSON))
	// http.Redirect(w, r, url.URL, 200)
}

func Update(w http.ResponseWriter, r *http.Request) {

}

func Delete(w http.ResponseWriter, r *http.Request) {

}

func GetStats(w http.ResponseWriter, r *http.Request) {

}
