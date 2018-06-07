package main

import (
	sh "github.com/murphysean/secrethitler"
	"time"
)

func GameFromGame(g sh.Game) Game {
	ret := Game{}
	ret.ID = g.ID
	ret.Secret = g.Secret
	ret.EventID = g.EventID
	ret.State = g.State
	ret.Draw = g.Draw
	ret.Discard = g.Discard
	ret.Liberal = g.Liberal
	ret.Facist = g.Facist
	ret.ElectionTracker = g.ElectionTracker
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
			Status:         p.Status,
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
	ret.PreviousPresidentID = g.PreviousPresidentID
	ret.PreviousChancellorID = g.PreviousChancellorID
	ret.PreviousEnactedPolicy = g.PreviousEnactedPolicy
	ret.SpecialElectionRoundID = g.SpecialElectionRoundID
	ret.SpecialElectionPresidentID = g.SpecialElectionPresidentID
	ret.WinningParty = g.WinningParty

	return ret
}

type Game struct {
	ID                         string       `json:"id"`
	Secret                     string       `json:"secret"`
	EventID                    int          `json:"eventId"`
	State                      string       `json:"state"`
	Draw                       []string     `json:"draw"`
	Discard                    []string     `json:"discard"`
	Liberal                    int          `json:"liberal"`
	Facist                     int          `json:"facist"`
	ElectionTracker            int          `json:"electionTracker"`
	Players                    []GamePlayer `json:"players"`
	Round                      Round        `json:"round"`
	NextPresidentID            string       `json:"nextPresidentId"`
	PreviousPresidentID        string       `json:"previousPresidentId"`
	PreviousChancellorID       string       `json:"previousChancellorId"`
	PreviousEnactedPolicy      string       `json:"previousEnactedPolicy"`
	SpecialElectionRoundID     int          `json:"specialElectionRoundId"`
	SpecialElectionPresidentID string       `json:"specialElectionPresidentId"`
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
	Status         string    `json:"status"`
}

type Round struct {
	ID              int      `json:"id"`
	PresidentID     string   `json:"presidentId"`
	ChancellorID    string   `json:"chancellorId"`
	State           string   `json:"state"`
	Votes           []Vote   `json:"votes"`
	Policies        []string `json:"policies"`
	EnactedPolicy   string   `json:"enactedPolicy"`
	ExecutiveAction string   `json:"executiveAction"`
}

type Vote struct {
	PlayerID string `json:"playerId"`
	Vote     bool   `json:"vote"`
}
