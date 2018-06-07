package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	sh "github.com/murphysean/secrethitler"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Writer struct {
	Name string
}

func (w Writer) Write(b []byte) (int, error) {
	f, err := os.OpenFile(w.Name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println("write:", w.Name, err)
		return 0, err
	}
	defer f.Close()

	return f.Write(b)
}

func (ah APIHandler) CreateGameHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	gameID := GenUUIDv4()
	game := sh.NewSecretHitler()
	game.ID = gameID
	game.Log = Writer{"games/" + gameID + ".json"}
	//Drop a game update event that sets the gameID
	actx := context.Background()
	actx = context.WithValue(actx, "playerID", "engine")
	err := game.SubmitEvent(actx, sh.GameEvent{
		BaseEvent: sh.BaseEvent{Type: sh.TypeGameUpdate},
		Game:      game.Game,
	})
	if err != nil {
		http.Error(w, JsonErrorString(err.Error()), http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	ah.m.Lock()
	ah.ActiveGames[gameID] = game
	ah.m.Unlock()

	e := json.NewEncoder(w)
	fg := GameFromGame(game.Game.Filter(r.Context()))
	e.Encode(&fg)
}

func (ah APIHandler) GetGamesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	type sg struct {
		ID      string `json:"id"`
		State   string `json:"state"`
		Players int    `json:"players"`
	}
	ret := make([]sg, 0)
	ah.m.RLock()
	for _, shg := range ah.ActiveGames {
		ret = append(ret, sg{
			ID:      shg.Game.ID,
			State:   shg.Game.State,
			Players: len(shg.Game.Players)})
	}
	ah.m.RUnlock()
	e := json.NewEncoder(w)
	e.Encode(&ret)
}

var gre = regexp.MustCompile(`^/api/games/([^/]+)/?.*$`)

func (ah APIHandler) GetGameHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rer := gre.FindStringSubmatch(r.URL.Path)
	if len(rer) != 2 {
		http.Error(w, JsonErrorString("No GameID found"), http.StatusBadRequest)
		return
	}
	var g sh.Game
	ah.m.RLock()
	ret, ok := ah.ActiveGames[rer[1]]
	ah.m.RUnlock()

	if !ok {
		http.Error(w, JsonErrorString("Not Found"), http.StatusNotFound)
		return
	}

	g = ret.Game
	e := json.NewEncoder(w)
	//Filter it for the authenticated user
	fg := GameFromGame(g.Filter(r.Context()))
	e.Encode(&fg)
}

func (ah APIHandler) UpdateGameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("playerID").(string) != "admin" {
		http.Error(w, JsonErrorString("Forbidden"), http.StatusForbidden)
		return
	}
	rer := gre.FindStringSubmatch(r.URL.Path)
	if len(rer) != 2 {
		http.Error(w, JsonErrorString("No GameID found"), http.StatusBadRequest)
		return
	}
	ah.m.RLock()
	ret, ok := ah.ActiveGames[rer[1]]
	ah.m.RUnlock()
	if !ok {
		http.Error(w, JsonErrorString("Not Found"), http.StatusNotFound)
		return
	}

	//Read in the game
	g := sh.Game{}
	d := json.NewDecoder(r.Body)
	err := d.Decode(&g)
	if err != nil {
		fmt.Println(err)
		http.Error(w, JsonErrorString(err.Error()), http.StatusBadRequest)
		return
	}
	g.ID = rer[1]

	//Set it
	err = ret.SubmitEvent(r.Context(), sh.GameEvent{
		BaseEvent: sh.BaseEvent{Type: sh.TypeGameUpdate},
		Game:      g,
	})

	//Return the new game
	e := json.NewEncoder(w)
	//Filter it for the authenticated user
	fg := GameFromGame(g.Filter(r.Context()))
	e.Encode(&fg)
}

func (ah APIHandler) CreateGameEventHandler(w http.ResponseWriter, r *http.Request) {
	rer := gre.FindStringSubmatch(r.URL.Path)
	if len(rer) != 2 {
		http.Error(w, JsonErrorString("No GameID found"), http.StatusBadRequest)
		return
	}
	ah.m.RLock()
	ret, ok := ah.ActiveGames[rer[1]]
	ah.m.RUnlock()
	if !ok {
		http.Error(w, JsonErrorString("Not Found"), http.StatusNotFound)
		return
	}
	//Read the whole body into a buffer (to be read twice)
	b, err := ioutil.ReadAll(r.Body)
	e, err := sh.UnmarshalEvent(b)
	if err != nil {
		http.Error(w, JsonErrorString(err.Error()), http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	if e == nil {
		http.Error(w, JsonErrorString("Nil Event"), http.StatusBadRequest)
		return
	}
	//Sanitize input
	switch e.GetType() {
	case sh.TypePlayerMessage:
		me := e.(sh.MessageEvent)
		me.Message = bluemonday.UGCPolicy().Sanitize(me.Message)
		e = me
	case sh.TypeReactEventID:
		fallthrough
	case sh.TypeReactPlayer:
		fallthrough
	case sh.TypeReactStatus:
		re := e.(sh.ReactEvent)
		re.Reaction = bluemonday.UGCPolicy().Sanitize(re.Reaction)
		e = re
	}
	//Validate & submit the event against the game state
	err = ret.SubmitEvent(r.Context(), e)
	if err != nil {
		http.Error(w, JsonErrorString(err.Error()), http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	enc := json.NewEncoder(w)
	enc.Encode(&e)
}

func shouldSendState(t string) bool {
	if strings.HasPrefix(t, "request.") {
		return false
	}
	if strings.HasPrefix(t, "react.") {
		return false
	}
	if t == sh.TypeGameInformation {
		return false
	}
	if t == sh.TypePlayerVote {
		return false
	}
	if t == sh.TypePlayerMessage {
		return false
	}
	if t == sh.TypeGuess {
		return false
	}
	return true
}

func (ah APIHandler) GetGameEventsHandler(w http.ResponseWriter, r *http.Request) {
	// https://www.html5rocks.com/en/tutorials/eventsource/basics/
	rer := gre.FindStringSubmatch(r.URL.Path)
	if len(rer) != 2 {
		http.Error(w, JsonErrorString("No GameID found"), http.StatusBadRequest)
		return
	}

	//If the game isn't in the active games list, it still might be a log file...
	ah.m.RLock()
	ret, ok := ah.ActiveGames[rer[1]]
	ah.m.RUnlock()

	leids := r.Header.Get("Last-Event-Id")
	leid, _ := strconv.Atoi(leids)
	//If the last event id is less than 0, or greater than a billion
	if leid < 0 || leid >= 1000000000 {
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, JsonErrorString("eof"), http.StatusTooManyRequests)
		return
	}

	if !ok {
		//Ensure that the game doesn't already exist
		if _, err := os.Stat("games/" + rer[1] + ".json"); err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, JsonErrorString("Not Found"), http.StatusNotFound)
			return
		}
	}
	//Is this a flushable connection
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "webserver doesn't support flushing", http.StatusInternalServerError)
		return
	}
	cnot, ok := w.(http.CloseNotifier)
	if !ok {
		http.Error(w, "webserver doesn't support closenotify", http.StatusInternalServerError)
		return
	}
	cnotchan := cnot.CloseNotify()

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, ": Getting Started\n\n")
	flusher.Flush()
	var err error
	if ret != nil {
		for _, p := range ret.Game.Players {
			var pb []byte
			if p, err := GetPlayer(r.Context(), p.ID); err != nil {
				p := new(Player)
				p.ID = p.ID
				p.Email = p.ID
				p.Name = p.ID
				p.Username = p.ID
				p.ThumbnailURL = "http://www.gravatar.com/avatar"
				pb, _ = json.Marshal(&p)
			} else {
				pb, _ = json.Marshal(&p)
			}
			fmt.Fprintf(w, "event: %s\n", "player")
			fmt.Fprintf(w, "data: %s\n\n", pb)
		}
		flusher.Flush()
	}

	myChan := make(chan sh.Event)
	geid := 0
	over := true
	if ret != nil {
		geid = ret.Game.EventID
		if ret.Game.WinningParty == "" {
			over = false
		}

	}
	if geid == 0 || leid < geid {
		//Stream events from the last event id specified
		go func() {
			f, err := os.OpenFile("games/"+rer[1]+".json", os.O_RDONLY, 0644)
			if err != nil {
				fmt.Println(err)
				return
			}
			err = sh.ReadEventLog(f, myChan)
			if err != nil && err != io.EOF {
				fmt.Println(err)
			}
		}()
		tg := sh.Game{}
		for e := range myChan {
			tg, _, err = tg.Apply(e)
			if err != nil {
				fmt.Println(err)
			}
			if leid > 0 && e.GetID() <= leid {
				continue
			}
			//Don't filter if the real game is over
			if !over {
				//Before sending an event, filter it for the auth'd user
				e = e.Filter(r.Context())
			}
			b, err := json.Marshal(&e)
			if err != nil {
				fmt.Println(err)
			}
			if e.GetType() == sh.TypePlayerJoin {
				var pb []byte
				pje := e.(sh.PlayerEvent)
				if p, err := GetPlayer(r.Context(), pje.Player.ID); err != nil {
					p := new(Player)
					p.ID = pje.Player.ID
					p.Email = pje.Player.ID
					p.Name = pje.Player.ID
					p.Username = pje.Player.ID
					p.ThumbnailURL = "http://www.gravatar.com/avatar"
					pb, _ = json.Marshal(&p)
				} else {
					pb, _ = json.Marshal(&p)
				}
				fmt.Fprintf(w, "id: %d\n", e.GetID())
				fmt.Fprintf(w, "event: %s\n", "player")
				fmt.Fprintf(w, "data: %s\n\n", pb)
			}
			fmt.Fprintf(w, "id: %d\n", e.GetID())
			fmt.Fprintf(w, "event: %s\n", e.GetType())
			fmt.Fprintf(w, "data: %s\n\n", b)
			//Only send state on mutating events
			if shouldSendState(e.GetType()) {
				fmt.Fprintf(w, "id: %d\n", e.GetID())
				fmt.Fprintf(w, "event: %s\n", "state")
				g := GameFromGame(tg)
				//Only filter if the real game is not over
				if !over {
					g = GameFromGame(tg.Filter(r.Context()))
				}
				b, err = json.Marshal(&g)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Fprintf(w, "data: %s\n\n", b)
				//Flush the data down the pipe
			}
			flusher.Flush()
		}
		//If the game is over, set last event id to a billion, return
		if over {
			fmt.Fprintf(w, "id: %d\n", 1000000000)
			fmt.Fprintf(w, "event: %s\n", "server.close")
			fmt.Fprintf(w, "data: %s\n\n", "{}")
		}
	}

	if ret == nil {
		flusher.Flush()
		return
	}

	myChan = make(chan sh.Event)

	//Subscribe to game events
	uid := GenUUIDv4()
	//Add this channel to the subscriber list for the game
	ret.AddSubscriber(uid, myChan)

	//Defer the removal of the chanel from the game on disconnect
	defer func() {
		ret.RemoveSubscriber(uid)
	}()

	//Loop on events coming out of the gameserver
	for {
		select {
		case e := <-myChan:
			if e == nil {
				flusher.Flush()
				return
			}
			if e.GetType() == sh.TypePlayerJoin {
				var pb []byte
				pje := e.(sh.PlayerEvent)
				if p, err := GetPlayer(r.Context(), pje.Player.ID); err != nil {
					p := new(Player)
					p.ID = pje.Player.ID
					p.Email = pje.Player.ID
					p.Name = pje.Player.ID
					p.Username = pje.Player.ID
					p.ThumbnailURL = "http://www.gravatar.com/avatar"
					pb, _ = json.Marshal(&p)
				} else {
					pb, _ = json.Marshal(&p)
				}
				fmt.Fprintf(w, "id: %d\n", e.GetID())
				fmt.Fprintf(w, "event: %s\n", "player")
				fmt.Fprintf(w, "data: %s\n\n", pb)
			}
			//Before sending an event, filter it for the auth'd user
			e = e.Filter(r.Context())
			b, err := json.Marshal(&e)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Fprintf(w, "id: %d\n", e.GetID())
			fmt.Fprintf(w, "event: %s\n", e.GetType())
			fmt.Fprintf(w, "data: %s\n\n", b)

			if shouldSendState(e.GetType()) {
				fmt.Fprintf(w, "id: %d\n", e.GetID())
				//Optionally also include a seperate event sending the whole state for the client to sync on
				fmt.Fprintf(w, "event: %s\n", "state")
				//Before sending the state, filter it for the auth'd user
				g := GameFromGame(ret.Game.Filter(r.Context()))
				b, _ := json.Marshal(&g)
				fmt.Fprintf(w, "data: %s\n\n", b)
			}
		case <-time.After(time.Minute):
			//Check the game state, if it is done, send the server close event
			if ret != nil && ret.State == sh.GameStateFinished {
				fmt.Fprintf(w, "id: %d\n", 1000000000)
				fmt.Fprintf(w, "event: %s\n", "server.close")
				fmt.Fprintf(w, "data: %s\n\n", "{}")
			} else {
				fmt.Fprintf(w, ": keepalive\n\n")
			}
		case <-cnotchan:
			flusher.Flush()
			return
		}
		//Flush the data down the pipe
		flusher.Flush()
	}
}
