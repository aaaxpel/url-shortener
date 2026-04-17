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

type ErrorResponse struct {
	Error string
}

func Create(w http.ResponseWriter, r *http.Request) {
	if r.PostFormValue("url") == "" {
		w.WriteHeader(400)
		response, _ := json.Marshal(ErrorResponse{Error: "No URL provided"})

		w.Write([]byte(response))
		return
	}

	length := 8                    // symbols
	byteLength := (length + 1) / 2 // hex encoding
	randomBytes := make([]byte, byteLength)
	rand.Read(randomBytes)

	shortCode := hex.EncodeToString(randomBytes)

	query := "INSERT INTO urls (short_code, url) VALUES (?, ?) RETURNING id, short_code, url, createdAt, updatedAt"
	newURL := Url{}
	// Doing it in 1 row doesn't guarantee that it's bad request or bad database query / server issue
	// Should ideally split it up to account for that
	err := db.QueryRow(query, shortCode, r.PostFormValue("url")).Scan(&newURL.ID, &newURL.ShortCode, &newURL.URL, &newURL.CreatedAt, &newURL.UpdatedAt)
	if err != nil {
		w.WriteHeader(400)
		log.Printf("error: %s", err)
		error, _ := json.Marshal(err)
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
	query := "SELECT id, short_code, url, createdAt, updatedAt FROM urls WHERE short_code = ?"
	url := Url{}
	err := db.QueryRow(query, r.PathValue("code")).Scan(&url.ID, &url.ShortCode, &url.URL, &url.CreatedAt, &url.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			response, _ := json.Marshal(ErrorResponse{Error: "URL not found with provided short code"})

			w.Write([]byte(response))
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

	// Optional redirect if we want to use this path as the shareable link
	// but probably best to have link/CODE instead for this
	// http.Redirect(w, r, url.URL, 200)
}

func Update(w http.ResponseWriter, r *http.Request) {
	if r.PostFormValue("url") == "" {
		w.WriteHeader(400)
		response, _ := json.Marshal(ErrorResponse{Error: "No URL provided"})

		w.Write([]byte(response))
		return
	}

	query := `UPDATE urls SET
		url = ?,
		updatedAt = CURRENT_TIMESTAMP
		WHERE short_code = ?
		RETURNING id, short_code, url, createdAt, updatedAt`
	url := Url{}
	err := db.QueryRow(query, r.PostFormValue("url"), r.PathValue("code")).Scan(&url.ID, &url.ShortCode, &url.URL, &url.CreatedAt, &url.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			response, _ := json.Marshal(ErrorResponse{Error: "URL not found with provided short code"})

			w.Write([]byte(response))
			return
		}
	}

	urlJSON, err := json.Marshal(url)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	w.WriteHeader(200)
	w.Write([]byte(urlJSON))
}

func Delete(w http.ResponseWriter, r *http.Request) {
	query := `DELETE FROM urls
		WHERE short_code = ?
		RETURNING id, short_code, url, createdAt, updatedAt`
	url := Url{}
	err := db.QueryRow(query, r.PostFormValue("url"), r.PathValue("code")).Scan(&url.ID, &url.ShortCode, &url.URL, &url.CreatedAt, &url.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			response, _ := json.Marshal(ErrorResponse{Error: "URL not found with provided short code"})

			w.Write([]byte(response))
			return
		}
	}

	urlJSON, err := json.Marshal(url)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	w.WriteHeader(204)
	w.Write([]byte(urlJSON))
}

func GetStats(w http.ResponseWriter, r *http.Request) {
	// Same code as func Get() but with accessCount returned as well
	query := "SELECT id, short_code, url, createdAt, updatedAt, accessCount FROM urls WHERE short_code = ?"
	url := Url{}
	err := db.QueryRow(query, r.PathValue("code")).Scan(&url.ID, &url.ShortCode, &url.URL, &url.CreatedAt, &url.UpdatedAt, &url.AccessCount)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
			return
		}
	}

	urlJSON, err := json.Marshal(url)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	w.WriteHeader(200)
	w.Write([]byte(urlJSON))
}
