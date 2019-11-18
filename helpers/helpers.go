package helpers

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"html/template"
	"log"
	"net/http"
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
