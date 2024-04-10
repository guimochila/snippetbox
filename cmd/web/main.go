package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// Serving static files
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// Http Handlers
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	log.Print("Server start on :4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
