package main

import (
	"context"
	"encoding/json"
	"fmt"
	sh "github.com/murphysean/secrethitler"
	"net/http"
	"os"
	"strconv"
	"time"
)

// https://www.html5rocks.com/en/tutorials/eventsource/basics/
func ServerSentEventsHandler(w http.ResponseWriter, r *http.Request) {
	//TODO Find room?

	//Is this a flushable connection
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "webserver doesn't support flushing", http.StatusInternalServerError)
		return
	}
	cnot, ok := w.(http.CloseNotifier)
	if !ok {
		http.Error(w, "webserver doesn't support close notifications", http.StatusInternalServerError)
		return
	}
	cnotchan := cnot.CloseNotify()
	myChan := make(chan sh.Event)

	//TODO Set this with the authenticated users playerID
	ctx := context.WithValue(r.Context(), "playerID", r.URL.Query().Get("playerID"))

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "*")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, ": Getting Started\n\n")
	flusher.Flush()

	leids := r.Header.Get("Last-Event-Id")
	leid, err := strconv.Atoi(leids)
	if err == nil {
		if leid < theGame.Game.EventID {
			//Stream events from the last event id specified
			go func() {
				f, err := os.OpenFile(theGameFile, os.O_RDONLY, 0644)
				if err != nil {
					fmt.Println(err)
					return
				}
				err = sh.ReadEventLog(f, myChan)
			}()
			tg := sh.Game{}
			for e := range myChan {
				tg, _, err = tg.Apply(e)
				if err != nil {
					fmt.Println(err)
				}
				if e.GetID() <= leid {
					continue
				}
				//TODO Don't filter if the real game is over
				//Before sending an event, filter it for the auth'd user
				e = e.Filter(ctx)
				b, err := json.Marshal(&e)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Fprintf(w, "id: %d\n", e.GetID())
				fmt.Fprintf(w, "event: %s\n", e.GetType())
				fmt.Fprintf(w, "data: %s\n\n", b)
				fmt.Fprintln(w, "event: state")
				//TODO Don't filter if the real game is over
				g := GameFromGame(tg.Filter(ctx))
				b, err = json.Marshal(&g)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Fprintf(w, "data: %s\n\n", b)
				//Flush the data down the pipe
				flusher.Flush()
			}
		}
	}
	myChan = make(chan sh.Event)

	//Subscribe to game events
	uid := GenUUIDv4()
	//Add this channel to the subscriber list for the game
	theGame.AddSubscriber(uid, myChan)

	//Defer the removal of the chanel from the game on disconnect
	defer func() {
		theGame.RemoveSubscriber(uid)
	}()

	//Loop on events coming out of the gameserver
	for {
		select {
		case e := <-myChan:
			if e == nil {
				flusher.Flush()
				return
			}
			//TODO Don't filter if the real game is over
			//Before sending an event, filter it for the auth'd user
			e = e.Filter(ctx)
			b, err := json.Marshal(&e)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Fprintf(w, "id: %d\n", e.GetID())
			fmt.Fprintf(w, "event: %s\n", e.GetType())
			fmt.Fprintf(w, "data: %s\n\n", b)
		case <-time.After(time.Minute):
			fmt.Fprintf(w, ": keepalive\n\n")
		case <-cnotchan:
			flusher.Flush()
			return
		}
		//Optionally also include a seperate event sending the whole state for the client to sync on
		fmt.Fprintln(w, "event: state")
		//TODO Don't filter if the real game is over
		//Before sending the state, filter it for the auth'd user
		g := GameFromGame(theGame.Game.Filter(ctx))
		b, err := json.Marshal(&g)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Fprintf(w, "data: %s\n\n", b)
		//Flush the data down the pipe
		flusher.Flush()
	}
}
