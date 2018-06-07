package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
)

func (ah APIHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Header.Get("Content-Type") {
	case "application/json":
		creds := struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}{}
		d := json.NewDecoder(r.Body)
		err := d.Decode(&creds)
		if err != nil {
			http.Error(w, JsonErrorString("Bad Request"), http.StatusBadRequest)
			return
		}
		if creds.Username == "" || creds.Password == "" {
			http.Error(w, JsonErrorString("Bad Request"), http.StatusBadRequest)
			return
		}
		g, err := GetGravatar(creds.Username)
		if err != nil {
			http.Error(w, JsonErrorString("Unauthorized"), http.StatusUnauthorized)
			return
		}
		p, err := GetPlayer(r.Context(), g.ID)
		if err != nil {
			http.Error(w, JsonErrorString("Unauthorized"), http.StatusUnauthorized)
			return
		}
		if !passwordHashmatches(creds.Username, creds.Password, p.PasswordHash) {
			http.Error(w, JsonErrorString("Unauthorized"), http.StatusUnauthorized)
			return
		}
		sessionID := GenUUIDv4()
		ah.Sessions[sessionID] = p
		ret := struct {
			Token  string  `json:"token"`
			Player *Player `json:"player"`
		}{
			Token:  sessionID,
			Player: p,
		}
		e := json.NewEncoder(w)
		e.Encode(&ret)
	case "application/x-www-form-urlencoded":
		username := r.PostFormValue("username")
		password := r.PostFormValue("password")
		if username == "" || password == "" {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		g, err := GetGravatar(username)
		if err != nil {
			http.Error(w, "Unauthorized"+err.Error(), http.StatusUnauthorized)
			return
		}
		p, err := GetPlayer(r.Context(), g.ID)
		if err != nil {
			http.Error(w, "Unauthorized"+err.Error(), http.StatusUnauthorized)
			return
		}
		if !passwordHashmatches(username, password, p.PasswordHash) {
			http.Error(w, "Unauthorized"+" password doesn't match", http.StatusUnauthorized)
			return
		}
		sessionID := GenUUIDv4()
		ah.Sessions[sessionID] = p
		//TODO Set secure?
		c := http.Cookie{
			HttpOnly: true,
			Name:     "shsid",
			Value:    sessionID,
			Path:     "/",
		}
		http.SetCookie(w, &c)
		//Redirect to index page
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	default:
		http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
		return
	}
}

func passwordHashmatches(username, password, hash string) bool {
	sig := hmac.New(sha256.New, []byte(username))
	sig.Write([]byte(password))
	passwordHash := base64.URLEncoding.EncodeToString(sig.Sum(nil))

	if passwordHash == hash {
		return true
	}

	return false
}
