package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	sh "github.com/murphysean/secrethitler"
	"net/http"
	"os"
	"strings"
	"sync"
)

var testingMode = true

func main() {
	//Specify a file to write all the events to

	os.MkdirAll("players", os.ModePerm)
	os.MkdirAll("games", os.ModePerm)

	apiHandler := NewAPIHandler()

	http.HandleFunc("/api/login", apiHandler.LoginHandler)
	http.Handle("/api/", apiHandler)

	//A file handler for the static assets
	http.Handle("/", http.FileServer(http.Dir("www/dist")))

	http.ListenAndServe(":8080", nil)
}

type APIHandler struct {
	Sessions    map[string]*Player
	ActiveGames map[string]*sh.SecretHitler
	m           sync.RWMutex
}

func NewAPIHandler() *APIHandler {
	ret := new(APIHandler)
	ret.ActiveGames = make(map[string]*sh.SecretHitler)
	ret.Sessions = make(map[string]*Player)
	return ret
}

func (ah *APIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//Set this with the authenticated users playerID
	ctx := r.Context()
	playerID := ""
	if h := r.Header.Get("Authorization"); h != "" && len(h) > 7 {
		token := h[7:]
		if ah.Sessions[token] != nil {
			playerID = ah.Sessions[token].ID
		}
	}
	if c, err := r.Cookie("shsid"); err == nil && playerID == "" {
		if ah.Sessions[c.Value] != nil {
			playerID = ah.Sessions[c.Value].ID
		}
	}
	if testingMode && playerID == "" {
		playerID = r.URL.Query().Get("playerID")
	}
	ctx = context.WithValue(r.Context(), "playerID", playerID)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "*")
	switch {
	case strings.HasPrefix(r.URL.Path, "/api/players"):
		if len(r.URL.Path) <= 13 {
			switch r.Method {
			case http.MethodPost:
				ah.CreatePlayerHandler(w, r.WithContext(ctx))
			default:
				http.Error(w, JsonErrorString("Method Not Allowed"), http.StatusMethodNotAllowed)
			}

		} else {
			//GET  /api/players/{playerID}
			ah.GetPlayerHandler(w, r.WithContext(ctx))
		}
	case strings.HasPrefix(r.URL.Path, "/api/games"):
		if len(r.URL.Path) <= 11 {
			//GET  /api/games
			//POST /api/games
			switch r.Method {
			case http.MethodGet:
				ah.GetGamesHandler(w, r.WithContext(ctx))
			case http.MethodPost:
				ah.CreateGameHandler(w, r.WithContext(ctx))
			default:
				http.Error(w, JsonErrorString("Method Not Allowed"), http.StatusMethodNotAllowed)
			}
		} else {
			if strings.HasSuffix(r.URL.Path, "/events") || strings.HasSuffix(r.URL.Path, "/events/") {
				switch r.Method {
				case http.MethodGet:
					//GET /api/games/{gameID}/events <- Get the event stream
					ah.GetGameEventsHandler(w, r.WithContext(ctx))
				case http.MethodPost:
					//POST /api/games/{gameID}/events <- Put a player event
					ah.CreateGameEventHandler(w, r.WithContext(ctx))
				}
			} else {
				switch r.Method {
				case http.MethodGet:
					//GET /api/games/{gameID}
					ah.GetGameHandler(w, r.WithContext(ctx))
				case http.MethodPut:
					//PUT /api/games/{gameID}
					ah.UpdateGameHandler(w, r.WithContext(ctx))
				}
			}
		}
	default:
		http.Error(w, JsonErrorString("Not Found"), http.StatusNotFound)
	}

}

func GenUUIDv4() string {
	u := make([]byte, 16)
	rand.Read(u)
	//Set the version to 4
	u[6] = (u[6] | 0x40) & 0x4F
	u[8] = (u[8] | 0x80) & 0xBF
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}
func JsonErrorString(s string) string {
	o := struct {
		Err string `json:"err"`
	}{s}
	b, _ := json.Marshal(&o)
	return string(b)
}
