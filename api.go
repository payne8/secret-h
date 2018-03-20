package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func APIStateHandler(w http.ResponseWriter, r *http.Request) {
	//Get the current game state
	w.Header().Set("Content-Type", "application/json")
	e := json.NewEncoder(w)
	//TODO Filter it for the authenticated user
	e.Encode(theGame.Game)
}

func APIEventHandler(w http.ResponseWriter, r *http.Request) {
	//Read the whole body into a buffer (to be read twice)
	b, err := ioutil.ReadAll(r.Body)
	//First determine the type of event being posted
	bt := BaseEvent{}
	err = json.Unmarshal(b, &bt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	//Second serialize it into that type of event object
	var e Event
	//Set the playerID on the event from authenticated context
	switch bt.GetType() {
	case TypePlayerJoin:
		fallthrough
	case TypePlayerReady:
		fallthrough
	case TypePlayerAcknowledge:
		pe := PlayerEvent{}
		err = json.Unmarshal(b, &pe)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println(err)
			return
		}
		e = pe
	case TypePlayerVote:
		pe := PlayerVoteEvent{}
		err = json.Unmarshal(b, &pe)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println(err)
			return
		}
		e = pe
	case TypePlayerNominate:
	case TypePlayerLegislate:
	case TypePlayerInvestigate:
	case TypePlayerSpecialElection:
	default:
		http.Error(w, "Unrecognized Event Type", http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	//Validate the event against the game state
	if e == nil {
		http.Error(w, "Nil Event", http.StatusBadRequest)
		return
	}
	err = theGame.Validate(r.Context(), e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	err = theGame.SubmitEvent(e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
