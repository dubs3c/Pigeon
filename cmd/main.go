package main

import (
	"log"
	"net/http"

	"github.com/mjdubell/Pigeon/pkg/onetimesecret"
)

const (
	staticDir = "web/static"
	port      = "9999"
)

func main() {

	db, _ := onetimesecret.NewDB()

	fs := http.FileServer(http.Dir(staticDir))

	router := onetimesecret.Router(db)
	defer db.Close()
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/", router)

	log.Println("Starting web server...")
	log.Fatal(http.ListenAndServe("127.0.0.1:"+port, nil))

}
