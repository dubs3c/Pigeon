package onetimesecret

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type data struct {
	Error   bool
	Message string
}

var repository *DB

// Router : HTTP Handler function
func Router(db *DB) *mux.Router {
	repository = db
	r := mux.NewRouter()

	r.HandleFunc("/", IndexGetHandler).Methods("GET")
	r.HandleFunc("/secret", IndexPostHandler).Methods("POST")
	r.HandleFunc("/secret/{token}", GetSecretHandler).Methods("GET")
	r.HandleFunc("/secret/{token}/unlock", GetPasswordProtectedSecretHandler).Methods("POST")

	return r
}

// IndexGetHandler : First page handler
func IndexGetHandler(w http.ResponseWriter, r *http.Request) {
	RenderTemplate("index", nil, w, r)
}

// GetSecretHandler : Fethes a secret by a given token
func GetSecretHandler(w http.ResponseWriter, r *http.Request) {
	token := mux.Vars(r)["token"]

	err := repository.IncrementViews(token)

	if err != nil {
		Return404(w, r)
		return
	}

	secret, err := repository.GetSecretByToken(token)

	if err != nil {
		Return404(w, r)
		return
	}

	currentTime := time.Now()
	valid := IsSecretValid(currentTime, secret)

	if !valid {
		repository.DeleteSecret(secret.token)
		Return404(w, r)
		return
	}

	var response data

	if secret.password != "" {
		response = data{
			Error:   true,
			Message: "This secret is password protected, please enter the correct password to unlock the secret.",
		}
	} else {
		response = data{
			Error:   false,
			Message: secret.message,
		}
	}

	RenderTemplate("secret", response, w, r)

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

	secret, err := repository.GetSecretByTokenAndPassword(token, userPassword)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "No secret found with that password"})
		return
	}

	currentTime := time.Now()
	valid := IsSecretValid(currentTime, secret)

	if !valid {
		repository.DeleteSecret(secret.token)
		Return404(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"secret": secret.message})

}

// IndexPostHandler : Create a secret
func IndexPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	maxview := r.PostForm.Get("maxview")

	if maxview == "" {
		maxview = "1"
	}
	maxviews, err := strconv.Atoi(maxview)

	if err != nil {
		log.Println("maxviews has to be an integer")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "maxviews is not an int!"})
		return
	}

	hash, err := GenerateToken()

	if err != nil {
		log.Print("Could not generate token for some reason: ", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Something went wrong :/"})
	}

	secret := Secret{
		token:    hash,
		message:  r.PostForm.Get("secret"),
		password: r.PostForm.Get("password"),
		expire:   r.PostForm.Get("expire"),
		maxviews: maxviews,
		views:    0,
	}

	err = repository.CreateSecret(secret)

	if err != nil {
		log.Println("Could not insert into database: ", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Opsie, you messed up"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": secret.token})
}
