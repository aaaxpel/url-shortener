package main

import (
	"net/http"
)

func Create(w http.ResponseWriter, r *http.Request) {

	// db := database.Connect()
	r.PostFormValue("url")

	w.WriteHeader(201)
	w.Write([]byte("hi"))
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
