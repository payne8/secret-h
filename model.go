package main

import (
	sh "github.com/murphysean/secrethitler"
	"time"
)

func GameFromGame(g sh.Game) Game {
	ret := Game{}
	ret.ID = g.ID
	ret.EventID = g.EventID
	ret.State = g.State
	ret.Draw = g.Draw
	ret.Discard = g.Discard
	ret.Liberal = g.Liberal
	ret.Facist = g.Facist
	ret.FailedVotes = g.FailedVotes
	ret.Players = []GamePlayer{}
	for _, p := range g.Players {
		np := GamePlayer{
			ID:             p.ID,
			Party:          p.Party,
			Role:           p.Role,
			Ready:          p.Ready,
			Ack:            p.Ack,
			ExecutedBy:     p.ExecutedBy,
			InvestigatedBy: p.InvestigatedBy,
			LastAction:     p.LastAction,
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
	ID                         string       `json:"id"`
	EventID                    int          `json:"eventID"`
	State                      string       `json:"state"`
	Draw                       []string     `json:"draw"`
	Discard                    []string     `json:"discard"`
	Liberal                    int          `json:"liberal"`
	Facist                     int          `json:"facist"`
	FailedVotes                int          `json:"failedVotes"`
	Players                    []GamePlayer `json:"players"`
	Round                      Round        `json:"round"`
	NextPresidentID            string       `json:"nextPresidentID"`
	PreviousPresidentID        string       `json:"previousPresidentID"`
	PreviousChancellorID       string       `json:"previousChancellorID"`
	SpecialElectionRoundID     int          `json:"specialElectionRoundID"`
	SpecialElectionPresidentID string       `json:"specialElectionPresidentID"`
	WinningParty               string       `json:"winningParty"`
}

type GamePlayer struct {
	ID             string    `json:"id"`
	Party          string    `json:"party"`
	Role           string    `json:"role"`
	Ready          bool      `json:"ready"`
	Ack            bool      `json:"ack"`
	ExecutedBy     string    `json:"executedBy"`
	InvestigatedBy string    `json:"investigatedBy"`
	LastAction     time.Time `json:"lastAction"`
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
