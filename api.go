package main

import (
	"encoding/json"
	"fmt"
	sh "github.com/murphysean/secrethitler"
	"io/ioutil"
	"net/http"
)

func APIStateHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		//Get the current game state
		w.Header().Set("Content-Type", "application/json")
		e := json.NewEncoder(w)
		//TODO Filter it for the authenticated user
		e.Encode(theGame.Game)
	case http.MethodPut:
		//TODO Only admins can do this
		g := sh.Game{}
		d := json.NewDecoder(r.Body)
		err := d.Decode(&g)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println(err)
		}
		//Fire a game update event with the changes
		err = theGame.SubmitEvent(r.Context(), sh.GameEvent{
			BaseEvent: sh.BaseEvent{Type: sh.TypeGameUpdate},
			Game:      g,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			fmt.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		e := json.NewEncoder(w)
		//TODO Filter it for the authenticated user
		e.Encode(theGame.Game)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func APIEventHandler(w http.ResponseWriter, r *http.Request) {
	//Read the whole body into a buffer (to be read twice)
	b, err := ioutil.ReadAll(r.Body)
	//First determine the type of event being posted
	bt := sh.BaseEvent{}
	err = json.Unmarshal(b, &bt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	//Second serialize it into that type of event object
	var e sh.Event
	//Set the playerID on the event from authenticated context
	switch bt.GetType() {
	case sh.TypePlayerJoin:
		fallthrough
	case sh.TypePlayerReady:
		fallthrough
	case sh.TypePlayerAcknowledge:
		pe := sh.PlayerEvent{}
		err = json.Unmarshal(b, &pe)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println(err)
			return
		}
		e = pe
	case sh.TypePlayerVote:
		pe := sh.PlayerVoteEvent{}
		err = json.Unmarshal(b, &pe)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println(err)
			return
		}
		e = pe
	case sh.TypePlayerNominate:
		pe := sh.PlayerPlayerEvent{}
		err = json.Unmarshal(b, &pe)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println(err)
			return
		}
		e = pe
	case sh.TypePlayerLegislate:
		pe := sh.PlayerLegislateEvent{}
		err = json.Unmarshal(b, &pe)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println(err)
			return
		}
		e = pe
	case sh.TypePlayerInvestigate:
		pe := sh.PlayerPlayerEvent{}
		err = json.Unmarshal(b, &pe)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println(err)
			return
		}
		e = pe
	case sh.TypePlayerSpecialElection:
		pe := sh.PlayerPlayerEvent{}
		err = json.Unmarshal(b, &pe)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println(err)
			return
		}
		e = pe
	case sh.TypePlayerExecute:
		pe := sh.PlayerPlayerEvent{}
		err = json.Unmarshal(b, &pe)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println(err)
			return
		}
		e = pe
	default:
		http.Error(w, "Unrecognized Event Type", http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	if e == nil {
		http.Error(w, "Nil Event", http.StatusBadRequest)
		return
	}
	//Validate & submit the event against the game state
	err = theGame.SubmitEvent(r.Context(), e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
