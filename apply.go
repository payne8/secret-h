package main

import (
	"time"
)

//Apply mutates the game state by applying the given event.
func (g Game) Apply(e Event) (Game, Event, error) {
	//Increment the event counter
	g.EventID = g.EventID + 1

	//Assign the event id to the event
	switch e.GetType() {
	case TypePlayerJoin:
		ne := e.(PlayerEvent)
		ne.ID = g.EventID
		ne.Moment = time.Now()
		e = ne
		g.Players = append(g.Players, ne.Player)
	case TypePlayerReady:
		ne := e.(PlayerEvent)
		ne.ID = g.EventID
		ne.Moment = time.Now()
		e = ne
		for i, p := range g.Players {
			if p.ID == ne.Player.ID {
				g.Players[i].Ready = true
				break
			}
		}
	case TypePlayerAcknowledge:
		ne := e.(PlayerEvent)
		ne.ID = g.EventID
		ne.Moment = time.Now()
		e = ne
		//Switch the given users ack attribute to true
		for i, p := range g.Players {
			if p.ID == ne.Player.ID {
				g.Players[i].Ack = true
				break
			}
		}
	case TypePlayerVote:
		ne := e.(PlayerVoteEvent)
		ne.ID = g.EventID
		ne.Moment = time.Now()
		e = ne
		//Add the given vote to the rounds vote array
		g.Round.Votes = append(g.Round.Votes, ne.Vote)
	case TypePlayerNominate:
		ne := e.(PlayerPlayerEvent)
		ne.ID = g.EventID
		ne.Moment = time.Now()
		e = ne
		//Add the chancelor to the round object
		g.Round.ChancellorID = ne.OtherPlayerID
	case TypePlayerLegislate:
		ne := e.(PlayerLegislateEvent)
		ne.ID = g.EventID
		ne.Moment = time.Now()
		e = ne
		//Add to the discard pile
		g.Discard = append(g.Discard, ne.Discard)
		//clean up the round array
		for i, v := range g.Round.Policies {
			if v == ne.Discard {
				g.Round.Policies = remove(g.Round.Policies, i)
				break
			}
		}
		//TODO Eventually account for veto
	case TypePlayerInvestigate:
		ne := e.(PlayerPlayerEvent)
		ne.ID = g.EventID
		ne.Moment = time.Now()
		e = ne
		for i, p := range g.Players {
			if p.ID == ne.OtherPlayerID {
				g.Players[i].InvestigatedBy = ne.PlayerID
			}
		}
	case TypePlayerSpecialElection:
		ne := e.(PlayerPlayerEvent)
		ne.ID = g.EventID
		ne.Moment = time.Now()
		e = ne
		g.SpecialElectionPresidentID = ne.OtherPlayerID
		g.SpecialElectionRoundID = g.Round.ID + 1
	case TypePlayerExecute:
		ne := e.(PlayerPlayerEvent)
		ne.ID = g.EventID
		ne.Moment = time.Now()
		e = ne
		for i, p := range g.Players {
			if p.ID == ne.OtherPlayerID {
				g.Players[i].ExecutedBy = ne.PlayerID
			}
		}
	case TypeGameStart:
		ne := e.(GameStartEvent)
		ne.ID = g.EventID
		ne.Moment = time.Now()
		e = ne
		//Using the event data, initialize the game state
		g.State = GameStateInit
		g.Draw = ne.Draw
		g.NextPresidentID = g.Players[ne.InitialPresidentIndex].ID
		for i, _ := range g.Players {
			g.Players[i].Role = ne.Roles[i]
			if g.Players[i].Role == RoleHitler {
				g.Players[i].Party = PartyFacist
			} else if g.Players[i].Role == RoleFacist {
				g.Players[i].Party = PartyFacist
			} else {
				g.Players[i].Party = PartyLiberal
			}
		}
	case TypeGameShuffle:
		ne := e.(GameShuffleEvent)
		ne.ID = g.EventID
		ne.Moment = time.Now()
		e = ne
		//The event data, set the discard and draw pile accordingly
		g.Draw = ne.Draw
		g.Discard = []string{}
	case TypeGameEnd:
		ne := e.(GameEndEvent)
		ne.ID = g.EventID
		ne.Moment = time.Now()
		e = ne
		//Set the game state attribute to finished
		g.State = GameStateFinished
	case TypeRoundStart:
		ne := e.(RoundStartEvent)
		ne.ID = g.EventID
		ne.Moment = time.Now()
		e = ne
		//Set the round object from the event data
		g.Round.ID = ne.RoundID
		g.Round.PresidentID = ne.PresidentID
		g.Round.ChancellorID = ""
		g.Round.State = RoundStateNominating
		g.Round.Votes = []Vote{}
		g.Round.Policies = []string{}
		g.Round.ExecutiveAction = ""
	case TypeRoundNominateRequest:
		ne := e.(RoundRequestEvent)
		ne.ID = g.EventID
		ne.Moment = time.Now()
		e = ne
	case TypeRoundVoteStart:
		ne := e.(RoundStateEvent)
		ne.ID = g.EventID
		ne.Moment = time.Now()
		e = ne
		//Set the round state to voting
		g.Round.State = ne.State
	case TypeRoundVoteEnd:
		ne := e.(RoundVoteEndEvent)
		ne.ID = g.EventID
		ne.Moment = time.Now()
		e = ne
		g.Board.FailedVotes = ne.FailedVoteCount
		g.Round.State = ne.RoundState
	case TypeRoundLegislateRequest:
		//TODO noop
	case TypeRoundLegislateEnact:
		//TODO Change the board state
		if g.Board.FailedVotes >= 3 {
			g.PreviousPresidentID = ""
			g.PreviousChancellorID = ""
		}
		//TODO If a policy was played with an available executive action, trigger that
		//TODO If not, end the round
		g.Board.FailedVotes = 0
	case TypeRoundExecutiveActionRequest:
		//TODO Noop
	case TypeRoundExecutiveActionEnact:
		ne := e.(RoundExecutiveActionEnact)
		ne.ID = g.EventID
		ne.Moment = time.Now()
		e = ne
		g.Round.State = RoundStateFinished
	case TypeRoundEnd:
		//Finalize the round, trigger the next round
		g.PreviousPresidentID = g.Round.PresidentID
		g.PreviousChancellorID = g.Round.ChancellorID
	}

	//TODO Move this to a top handler, allowing us to indpeendently test application of an event to the state Broadcast the event
	return g, e, nil
}

func remove(s []string, i int) []string {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}
