package main

import (
	"log"
	"strconv"
	"net/http"
	"crypto/rand"
	"encoding/hex"
	"database/sql"
	"html/template"
	"encoding/json"
	"crypto/sha256"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Secret struct {
	token    string
	secret   string
	password string
	expire   string
	maxviews int
}

type Data struct {
	Error	bool
	Message string
}

const (
	STATIC_DIR = "/static/"
	PORT       = "9999"
)

var db *sql.DB

func renderTemplate(temp string, data interface{}, w http.ResponseWriter, r *http.Request) {
	files := []string{
		"templates/layout.html",
		"templates/" + temp + ".html",
	}
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatal("ERROR: ", err)
	}
	tmpl.ExecuteTemplate(w, "base", data)
}

func indexGetHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate("index", nil, w, r)
}

func generateToken() (string, error) {
	c := 32
	b := make([]byte, c)
	_, err := rand.Read(b)
	if err != nil {
		log.Println("Error generating random token: ", err)
		return "", err
	}
	h := sha256.New()
	h.Write(b)
	encodedToken := hex.EncodeToString(h.Sum(nil))[:32]
	return encodedToken, nil
}

func returnJson(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"data": "lol"})
}

func getSecretHandler(w http.ResponseWriter, r *http.Request) {
	token := mux.Vars(r)["token"]
	row := db.QueryRow("SELECT * FROM Secrets WHERE token=?", token)
	
	var (
		secret   string
		password string
		expire   string
		maxview  int
	)

	if err := row.Scan(&token, &secret, &password, &expire, &maxview); err != nil {
		log.Println(err)
		return
	}

	if password != "" {

		data := Data{
			Error: true,
			Message: "This secret is password protected, please enter the correct password to unlock the secret.",
		}
		renderTemplate("secret", data, w, r)

	} else {
		data := Data{
			Error: false,
			Message: secret,
		}
		renderTemplate("secret", data, w, r)
	}

}

func getPasswordProtectedSecretHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	token := mux.Vars(r)["token"]
	userPassword := r.PostForm.Get("password")

	if userPassword == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": "You specified an empty password"})
		return
	}

	row := db.QueryRow("SELECT * FROM Secrets WHERE token=? AND password=?", token, userPassword)
	
	var (
		secret   string
		password string
		expire   string
		maxview  int
	)

	if err := row.Scan(&token, &secret, &password, &expire, &maxview); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "No secret found with that password"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"secret": secret})

}

func indexPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	maxviews, err := strconv.Atoi(r.PostForm.Get("maxview"))

	if err != nil {
		log.Println("maxviews has to be an integer")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": "maxviews is not an int!"})
		return
	}


	hash, err := generateToken()

	if err != nil {
		log.Print("Could not generate token for some reason: ", err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": "Something went wrong :/"})
	}

	message := Secret{
		token:    hash,
		secret:   r.PostForm.Get("secret"),
		password: r.PostForm.Get("password"),
		expire:   r.PostForm.Get("expire"),
		maxviews: maxviews,
	}
	
	log.Println(message)

	stmt, err := db.Prepare("INSERT INTO Secrets VALUES(?, ?, ?, ?, ?)")
	
	if err != nil {
		log.Println("Could not insert into database: ", err)
		return
	}
	
	stmt.Exec(message.token, message.secret, message.password, "", message.maxviews)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": message.token})
}

func initDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "db.sqlite3")
	if err != nil {
		log.Fatal(err)
	}

	create := `
		CREATE TABLE IF NOT EXISTS Secrets (
		token VARCHAR(32) PRIMARY KEY,
		secret TEXT,
		password VARCHAR(255),
		expire TEXT,
		maxviews INTEGER DEFAULT 1
		);`

	stmt, err := db.Prepare(create)
	if err != nil {
		log.Fatal("Error creating table: ", err)
	}

	stmt.Exec()

	err = db.Ping()
	if err != nil {
		log.Fatal("Pinging database failed: ", err)
	}

	return db
}

func main() {

	db = initDatabase()

	r := mux.NewRouter()
	r.HandleFunc("/", indexGetHandler).Methods("GET")
	r.HandleFunc("/secret", indexPostHandler).Methods("POST")
	r.HandleFunc("/secret/{token}", getSecretHandler).Methods("GET")
	r.HandleFunc("/secret/{token}/unlock", getPasswordProtectedSecretHandler).Methods("POST")
	r.PathPrefix(STATIC_DIR).Handler(http.StripPrefix(STATIC_DIR, http.FileServer(http.Dir("."+STATIC_DIR))))
	http.Handle("/", r)
	log.Println("Starting web server...")
	log.Fatal(http.ListenAndServe(":"+PORT, nil))

	defer db.Close()
}
