package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/mjdubell/pigeon/helpers"
)

type secret struct {
	token    string
	secret   string
	password string
	expire   string
	maxviews int
	views    int
}

type data struct {
	Error   bool
	Message string
}

// DB : Database pointer containing the DB connection, set in main.go.
var DB *sql.DB

// IndexGetHandler : First page handler
func IndexGetHandler(w http.ResponseWriter, r *http.Request) {
	helpers.RenderTemplate("index", nil, w, r)
}

// GetSecretHandler : Fethes a secret by a given token
func GetSecretHandler(w http.ResponseWriter, r *http.Request) {
	token := mux.Vars(r)["token"]
	row := DB.QueryRow("SELECT * FROM Secrets WHERE token=?", token)

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

		data := data{
			Error:   true,
			Message: "This secret is password protected, please enter the correct password to unlock the secret.",
		}
		helpers.RenderTemplate("secret", data, w, r)

	} else {
		data := data{
			Error:   false,
			Message: secret,
		}
		helpers.RenderTemplate("secret", data, w, r)
	}

}

// GetPasswordProtectedSecretHandler : Fetches a password protected secret
func GetPasswordProtectedSecretHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	token := mux.Vars(r)["token"]
	userPassword := r.PostForm.Get("password")

	if userPassword == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": "You specified an empty password"})
		return
	}

	row := DB.QueryRow("SELECT * FROM Secrets WHERE token=? AND password=?", token, userPassword)

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

// IndexPostHandler : Create a secret
func IndexPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	maxviews, err := strconv.Atoi(r.PostForm.Get("maxview"))

	if err != nil {
		log.Println("maxviews has to be an integer")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": "maxviews is not an int!"})
		return
	}

	hash, err := helpers.GenerateToken()

	if err != nil {
		log.Print("Could not generate token for some reason: ", err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"error": "Something went wrong :/"})
	}

	message := secret{
		token:    hash,
		secret:   r.PostForm.Get("secret"),
		password: r.PostForm.Get("password"),
		expire:   r.PostForm.Get("expire"),
		maxviews: maxviews,
	}

	log.Println(message)

	stmt, err := DB.Prepare("INSERT INTO Secrets VALUES(?, ?, ?, ?, ?)")

	if err != nil {
		log.Println("Could not insert into database: ", err)
		return
	}

	stmt.Exec(message.token, message.secret, message.password, "", message.maxviews)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": message.token})
}
