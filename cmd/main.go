package main

import (
	"log"
	"net/http"

	"github.com/mjdubell/pigeon/handlers"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

const (
	staticDir = "web/static"
	port      = "9999"
)

func main() {

	db := InitDatabase()

	handlers.DB = db

	r := mux.NewRouter()
	fs := http.FileServer(http.Dir(staticDir))

	r.HandleFunc("/", handlers.IndexGetHandler).Methods("GET")
	r.HandleFunc("/secret", handlers.IndexPostHandler).Methods("POST")
	r.HandleFunc("/secret/{token}", handlers.GetSecretHandler).Methods("GET")
	r.HandleFunc("/secret/{token}/unlock", handlers.GetPasswordProtectedSecretHandler).Methods("POST")

	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/", r)

	log.Println("Starting web server...")
	log.Fatal(http.ListenAndServe(":"+port, nil))

	defer db.Close()
}
