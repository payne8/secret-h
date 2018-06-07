package main

import (
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Player struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	Name         string `json:"name"`
	ThumbnailURL string `json:"thumbnailUrl"`
	Password     string `json:"password,omitempty"`
	PasswordHash string `json:"passwordHash,omitempty"`
}

type Gravatar struct {
	ID                string `json:"id"`
	PreferredUseranme string `json:"preferredUseranme"`
	ThumbnailURL      string `json:"thumbnailUrl"`
	/*
		Name              struct {
			GivenName  string `json:"givenName"`
			FamilyName string `json:"familyName"`
			Formatted  string `json:"formatted"`
		} `json:"name"`
	*/
	DisplayName     string `json:"displayName"`
	CurrentLocation string `json:"currentLocation"`
}

func GetGravatar(email string) (*Gravatar, error) {
	//Look up the gravatar information
	te := strings.TrimSpace(email)
	le := strings.ToLower(te)
	h := md5.New()
	h.Write([]byte(le))
	hash := h.Sum(nil)
	ghash := fmt.Sprintf("%x", hash)

	resp, err := http.Get("https://www.gravatar.com/" + ghash + ".json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status + " response")
	}
	d := json.NewDecoder(resp.Body)
	ro := struct {
		Entry []Gravatar `json:"entry"`
	}{}
	err = d.Decode(&ro)
	if err != nil {
		return nil, err
	}

	if ro.Entry == nil && len(ro.Entry) < 1 {
		return nil, errors.New("Gravatar resp empty")
	}

	return &ro.Entry[0], nil
}

func GetPlayer(ctx context.Context, id string) (*Player, error) {
	if id == "me" {
		id, _ = ctx.Value("playerID").(string)
	}
	b, err := ioutil.ReadFile("players/" + id + ".json")
	if err != nil {
		return nil, err
	}
	p := Player{}
	err = json.Unmarshal(b, &p)
	if err != nil {
		return nil, err
	}

	return &p, err
}

func (ah APIHandler) CreatePlayerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//Given a email and a password, create this user

	//Pull in the player object
	p := Player{}
	d := json.NewDecoder(r.Body)
	err := d.Decode(&p)
	if err != nil {
		fmt.Println(err)
		http.Error(w, JsonErrorString(err.Error()), http.StatusBadRequest)
		return
	}

	g, err := GetGravatar(p.Email)
	if err != nil {
		fmt.Println(err)
		http.Error(w, JsonErrorString(err.Error()), http.StatusInternalServerError)
		return
	}

	//Ensure that the user doesn't already exist
	if _, err := os.Stat("players/" + p.ID + ".json"); err == nil {
		http.Error(w, JsonErrorString("User already exists"), http.StatusConflict)
		return
	}

	p.ID = g.ID
	p.Name = g.DisplayName
	p.ThumbnailURL = g.ThumbnailURL
	p.Username = g.DisplayName

	sig := hmac.New(sha256.New, []byte(p.Email))
	sig.Write([]byte(p.Password))
	p.Password = ""
	p.PasswordHash = base64.URLEncoding.EncodeToString(sig.Sum(nil))

	b, err := json.Marshal(&p)
	if err != nil {
		fmt.Println(err)
		http.Error(w, JsonErrorString(err.Error()), http.StatusInternalServerError)
		return
	}

	err = ioutil.WriteFile("players/"+p.ID+".json", b, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		http.Error(w, JsonErrorString(err.Error()), http.StatusInternalServerError)
		return
	}
	p.PasswordHash = ""

	w.Header().Set("Location", "/api/players/"+p.ID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	e := json.NewEncoder(w)
	e.Encode(&p)
}

var pre = regexp.MustCompile(`^/api/players/([^/]+)/?$`)

func (ah APIHandler) GetPlayerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rer := pre.FindStringSubmatch(r.URL.Path)
	if len(rer) != 2 {
		http.Error(w, JsonErrorString("No PlayerID found"), http.StatusBadRequest)
		return
	}

	p, err := GetPlayer(r.Context(), rer[1])
	if err != nil {
		http.Error(w, JsonErrorString("Not Found"), http.StatusNotFound)
		return
	}
	p.Password = ""
	p.PasswordHash = ""

	e := json.NewEncoder(w)
	e.Encode(&p)
}
