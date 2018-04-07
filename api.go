package main

import (
	"context"
	"encoding/json"
	"fmt"
	sh "github.com/murphysean/secrethitler"
	"io/ioutil"
	"net/http"
)

func GameFromGame(g sh.Game) Game {
	ret := Game{}
	ret.EventID = g.EventID
	ret.State = g.State
	ret.Draw = g.Draw
	ret.Discard = g.Discard
	ret.Liberal = g.Liberal
	ret.Facist = g.Facist
	ret.FailedVotes = g.FailedVotes
	ret.Players = []Player{}
	for _, p := range g.Players {
		np := Player{
			ID:             p.ID,
			Name:           p.Name,
			Party:          p.Party,
			Role:           p.Role,
			Ready:          p.Ready,
			Ack:            p.Ack,
			ExecutedBy:     p.ExecutedBy,
			InvestigatedBy: p.InvestigatedBy,
		}
		ret.Players = append(ret.Players, np)
	}
	ret.Round = Round{}
	ret.Round.ID = g.Round.ID
	ret.Round.PresidentID = g.Round.PresidentID
	ret.Round.ChancellorID = g.Round.ChancellorID
	ret.Round.State = g.Round.State
	ret.Round.Votes = []Vote{}
	for _, v := range g.Round.Votes {
		nv := Vote{PlayerID: v.PlayerID, Vote: v.Vote}
		ret.Round.Votes = append(ret.Round.Votes, nv)
	}
	ret.Round.Policies = g.Round.Policies
	ret.Round.EnactedPolicy = g.Round.EnactedPolicy
	ret.Round.ExecutiveAction = g.Round.ExecutiveAction
	ret.NextPresidentID = g.NextPresidentID
	ret.PreviousChancellorID = g.PreviousChancellorID
	ret.SpecialElectionRoundID = g.SpecialElectionRoundID
	ret.SpecialElectionPresidentID = g.SpecialElectionPresidentID
	ret.WinningParty = g.WinningParty

	return ret
}

type Game struct {
	EventID                    int      `json:"eventID"`
	State                      string   `json:"state"`
	Draw                       []string `json:"draw"`
	Discard                    []string `json:"discard"`
	Liberal                    int      `json:"liberal"`
	Facist                     int      `json:"facist"`
	FailedVotes                int      `json:"failedVotes"`
	Players                    []Player `json:"players"`
	Round                      Round    `json:"round"`
	NextPresidentID            string   `json:"nextPresidentID"`
	PreviousPresidentID        string   `json:"previousPresidentID"`
	PreviousChancellorID       string   `json:"previousChancellorID"`
	SpecialElectionRoundID     int      `json:"specialElectionRoundID"`
	SpecialElectionPresidentID string   `json:"specialElectionPresidentID"`
	WinningParty               string   `json:"winningParty"`
}

type Player struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Party          string `json:"party"`
	Role           string `json:"role"`
	Ready          bool   `json:"ready"`
	Ack            bool   `json:"ack"`
	ExecutedBy     string `json:"executedBy"`
	InvestigatedBy string `json:"investigatedBy"`
}

type Round struct {
	ID              int      `json:"id"`
	PresidentID     string   `json:"presidentID"`
	ChancellorID    string   `json:"chancellorID"`
	State           string   `json:"state"`
	Votes           []Vote   `json:"votes"`
	Policies        []string `json:"policies"`
	EnactedPolicy   string   `json:"enactedPolicy"`
	ExecutiveAction string   `json:"executiveAction"`
}

type Vote struct {
	PlayerID string `json:"playerID"`
	Vote     bool   `json:"vote"`
}

func APIStateHandler(w http.ResponseWriter, r *http.Request) {
	//TODO Set this with the authenticated users playerID
	ctx := context.WithValue(r.Context(), "playerID", "")
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		//Get the current game state
		e := json.NewEncoder(w)
		//TODO Filter it for the authenticated user
		e.Encode(GameFromGame(theGame.Game))
	case http.MethodPut:
		//TODO Only admins can do this
		g := sh.Game{}
		d := json.NewDecoder(r.Body)
		err := d.Decode(&g)
		if err != nil {
			http.Error(w, JsonErrorString(err.Error()), http.StatusBadRequest)
			fmt.Println(err)
		}
		//Fire a game update event with the changes
		err = theGame.SubmitEvent(ctx, sh.GameEvent{
			BaseEvent: sh.BaseEvent{Type: sh.TypeGameUpdate},
			Game:      g,
		})
		if err != nil {
			http.Error(w, JsonErrorString(err.Error()), http.StatusInternalServerError)
			fmt.Println(err)
		}
		e := json.NewEncoder(w)
		//Filter it for the authenticated user
		fg := GameFromGame(theGame.Game.Filter(ctx))
		e.Encode(&fg)
	default:
		http.Error(w, JsonErrorString("Method Not Allowed"), http.StatusMethodNotAllowed)
	}
}

func APIEventHandler(w http.ResponseWriter, r *http.Request) {
	//TODO Set this with the authenticated users playerID
	ctx := context.WithValue(r.Context(), "playerID", "")
	//Read the whole body into a buffer (to be read twice)
	b, err := ioutil.ReadAll(r.Body)
	//First determine the type of event being posted
	bt := sh.BaseEvent{}
	err = json.Unmarshal(b, &bt)
	if err != nil {
		http.Error(w, JsonErrorString(err.Error()), http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
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
			http.Error(w, JsonErrorString(err.Error()), http.StatusBadRequest)
			fmt.Println(err)
			return
		}
		e = pe
	case sh.TypePlayerVote:
		pe := sh.PlayerVoteEvent{}
		err = json.Unmarshal(b, &pe)
		if err != nil {
			http.Error(w, JsonErrorString(err.Error()), http.StatusBadRequest)
			fmt.Println(err)
			return
		}
		e = pe
	case sh.TypePlayerNominate:
		pe := sh.PlayerPlayerEvent{}
		err = json.Unmarshal(b, &pe)
		if err != nil {
			http.Error(w, JsonErrorString(err.Error()), http.StatusBadRequest)
			fmt.Println(err)
			return
		}
		e = pe
	case sh.TypePlayerLegislate:
		pe := sh.PlayerLegislateEvent{}
		err = json.Unmarshal(b, &pe)
		if err != nil {
			http.Error(w, JsonErrorString(err.Error()), http.StatusBadRequest)
			fmt.Println(err)
			return
		}
		e = pe
	case sh.TypePlayerInvestigate:
		pe := sh.PlayerPlayerEvent{}
		err = json.Unmarshal(b, &pe)
		if err != nil {
			http.Error(w, JsonErrorString(err.Error()), http.StatusBadRequest)
			fmt.Println(err)
			return
		}
		e = pe
	case sh.TypePlayerSpecialElection:
		pe := sh.PlayerPlayerEvent{}
		err = json.Unmarshal(b, &pe)
		if err != nil {
			http.Error(w, JsonErrorString(err.Error()), http.StatusBadRequest)
			fmt.Println(err)
			return
		}
		e = pe
	case sh.TypePlayerExecute:
		pe := sh.PlayerPlayerEvent{}
		err = json.Unmarshal(b, &pe)
		if err != nil {
			http.Error(w, JsonErrorString(err.Error()), http.StatusBadRequest)
			fmt.Println(err)
			return
		}
		e = pe
	default:
		http.Error(w, JsonErrorString("Unrecognized Event Type"), http.StatusBadRequest)
		fmt.Println(err)
		return
	}
	if e == nil {
		http.Error(w, JsonErrorString("Nil Event"), http.StatusBadRequest)
		return
	}
	//Validate & submit the event against the game state
	err = theGame.SubmitEvent(ctx, e)
	if err != nil {
		http.Error(w, JsonErrorString(err.Error()), http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	enc := json.NewEncoder(w)
	err = enc.Encode(&e)
	if err != nil {
		http.Error(w, JsonErrorString(err.Error()), http.StatusBadRequest)
		fmt.Println(err)
		return
	}

}

func JsonErrorString(s string) string {
	o := struct {
		Err string `json:"err"`
	}{s}
	b, _ := json.Marshal(&o)
	return string(b)
}
