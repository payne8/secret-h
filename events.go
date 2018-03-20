package main

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

const (
	GameStateLobby    = ""
	GameStateInit     = "init"
	GameStateStarted  = "started"
	GameStateFinished = "finished"

	RoundStateNominating      = "nominating"
	RoundStateVoting          = "voting"
	RoundStateFailed          = "failed"
	RoundStateLegislating     = "legislating"
	RoundStateExecutiveAction = "executive_action"
	RoundStateFinished        = "finished"

	ExecutiveActionInvestigate     = "investigate"
	ExecutiveActionPeek            = "peek"
	ExecutiveActionSpecialElection = "special_election"
	ExecutiveActionExecute         = "execute"

	TypePlayerJoin            = "player.join"
	TypePlayerReady           = "player.ready"
	TypePlayerAcknowledge     = "player.acknowledge"
	TypePlayerNominate        = "player.nominate"
	TypePlayerVote            = "player.vote"
	TypePlayerLegislate       = "player.legislate"
	TypePlayerInvestigate     = "player.investigate"
	TypePlayerSpecialElection = "player.special_election"
	TypePlayerExecute         = "player.execute"

	TypeGameStart   = "game.start"
	TypeGameShuffle = "game.shuffle"
	TypeGameEnd     = "game.end"

	TypeRoundStart                  = "round.start"
	TypeRoundNominateRequest        = "round.nominate_request"
	TypeRoundVoteStart              = "round.vote_start"
	TypeRoundVoteEnd                = "round.vote_end"
	TypeRoundLegislateRequest       = "round.legislate_request"
	TypeRoundLegislateEnact         = "round.legislate_enact"
	TypeRoundExecutiveActionRequest = "round.executive_action_request"
	TypeRoundExecutiveActionEnact   = "round.executive_action_enact"
	TypeRoundEnd                    = "round.end"
)

type Event interface {
	GetID() int
	GetType() string
}

type BaseEvent struct {
	ID     int       `json:"id"`
	Type   string    `json:"type"`
	Moment time.Time `json:"moment"`
}

func (e BaseEvent) GetID() int      { return e.ID }
func (e BaseEvent) GetType() string { return e.Type }

type PlayerEvent struct {
	BaseEvent
	Player Player `json:"player"`
}

type PlayerPlayerEvent struct {
	BaseEvent
	PlayerID      string `json:"playerID"`
	OtherPlayerID string `json:"otherPlayerID"`
}

type PlayerVoteEvent struct {
	BaseEvent
	Vote Vote `json:"vote"`
}

type PlayerLegislateEvent struct {
	BaseEvent
	PlayerID string
	Discard  string
	Veto     bool
}

func NewGameStartEvent(numPlayers int) GameStartEvent {
	ret := GameStartEvent{}
	ret.Type = TypeGameStart
	ret.Draw = make([]string, 0)
	for i := 0; i < 11; i++ {
		ret.Draw = append(ret.Draw, PolicyFacist)
	}
	for i := 0; i < 6; i++ {
		ret.Draw = append(ret.Draw, PolicyLiberal)
	}
	rand.Shuffle(len(ret.Draw), func(i, j int) {
		ret.Draw[i], ret.Draw[j] = ret.Draw[j], ret.Draw[i]
	})
	ret.Roles = []string{RoleLiberal, RoleLiberal, RoleLiberal, RoleHitler, RoleFacist}
	if numPlayers > 5 {
		ret.Roles = append(ret.Roles, RoleLiberal)
	}
	if numPlayers > 6 {
		ret.Roles = append(ret.Roles, RoleFacist)
	}
	if numPlayers > 7 {
		ret.Roles = append(ret.Roles, RoleLiberal)
	}
	if numPlayers > 8 {
		ret.Roles = append(ret.Roles, RoleFacist)
	}
	if numPlayers > 9 {
		ret.Roles = append(ret.Roles, RoleLiberal)
	}
	rand.Shuffle(len(ret.Roles), func(i, j int) {
		ret.Roles[i], ret.Roles[j] = ret.Roles[j], ret.Roles[i]
	})
	ret.InitialPresidentIndex = rand.Intn(numPlayers - 1)
	return ret
}

type GameStartEvent struct {
	BaseEvent
	Draw                  []string `json:"draw"`
	Roles                 []string `json:"roles"`
	InitialPresidentIndex int      `json:"initialPresidentIndex"`
}

type GameShuffleEvent struct {
	BaseEvent
	Draw []string `json:"draw"`
}

type GameEndEvent struct {
	BaseEvent
	WinningParty string `json:"winningParty"`
}

type RoundStartEvent struct {
	BaseEvent
	RoundID         int
	PresidentID     string
	NextPresidentID string
}

type RoundStateEvent struct {
	BaseEvent
	RoundID int
	State   string
}

type RoundRequestEvent struct {
	BaseEvent
	PlayerID int      `json:"playerID"`
	RoundID  int      `json:"roundID"`
	Policies []string `json:"policies,omitempty"`
}

type RoundLegislateEnact struct {
	BaseEvent
	RoundID int
	Policy  string
}

type RoundExecutiveActionEnact struct {
	BaseEvent
	RoundID           int
	ExecutiveAction   string
	PlayerID          string
	OtherPlayerID     string
	InvestigatedParty string
	PeekCards         []string
}

type RoundEndEvent struct {
	PreviousPresidentID  string `json:"previousPresidentID"`
	PreviousChancellorID string `json:"previousChancellorID"`
}

type RoundVoteEndEvent struct {
	BaseEvent
	RoundID         int
	Votes           []Vote
	Succeeded       bool
	FailedVoteCount int
	RoundState      string
}
