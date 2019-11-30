package onetimesecret

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"html/template"
	"log"
	"net/http"
	"time"
)

// RenderTemplate : Dynamically render templates
func RenderTemplate(temp string, data interface{}, w http.ResponseWriter, r *http.Request) {
	files := []string{
		"web/templates/layout.html",
		"web/templates/" + temp + ".html",
	}
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Fatal("ERROR: ", err)
	}
	tmpl.ExecuteTemplate(w, "base", data)
}

// GenerateToken : Generates a unique token for a secret
func GenerateToken() (string, error) {
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

// IsSecretValid : Check weather a given secret is valid
func IsSecretValid(currentTime time.Time, secret *Secret) bool {
	zone, offset := currentTime.Local().Zone()
	loc := time.FixedZone(zone, offset)

	expiredDatetime, err := time.ParseInLocation("2006-01-02 15:04:05", secret.expire, loc)
	if err != nil {
		log.Println("Error: Converting time failed: ", err)
		return false
	}
	return secret.views <= secret.maxviews && expiredDatetime.After(currentTime)
}

// Return404 : return page not found
func Return404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	RenderTemplate("secret", nil, w, r)
}
