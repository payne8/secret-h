package main

import ()

func nextIndex(len, idx int) int {
	if idx+1 >= len {
		return 0
	}
	return idx + 1
}

func (gs Game) createNextRound() RoundStartEvent {
	rse := RoundStartEvent{}
	rse.Type = TypeRoundStart
	rse.RoundID = gs.Round.ID + 1

	//Is the next round a special election?
	if rse.RoundID == gs.SpecialElectionRoundID {
		rse.PresidentID = gs.SpecialElectionPresidentID
		rse.NextPresidentID = gs.NextPresidentID
	} else {
		//Go to the next unexecuted president in the array
		//Next president is the next one in the array, that's not dead
		pi := -1
		for i, p := range gs.Players {
			if gs.NextPresidentID == p.ID {
				pi = i
				break
			}
		}
		for {
			//If president index is not dead, break
			if gs.Players[pi].ExecutedBy == "" {
				break
			} else {
				pi = nextIndex(len(gs.Players), pi)
			}
		}
		rse.PresidentID = gs.Players[pi].ID
		npi := nextIndex(len(gs.Players), pi)
		for {
			//If president index is not dead, break
			if gs.Players[npi].ExecutedBy == "" {
				break
			} else {
				npi = nextIndex(len(gs.Players), npi)
			}
		}
		rse.NextPresidentID = gs.Players[npi].ID
	}

	return rse
}

//TODO The engine will read the incoming event and process it to see if a new event
// should be created to update the game state. This function itself should not modify the game
// state in any way other than returning events that will.
func (g Game) Engine(e Event) ([]Event, error) {
	ret := []Event{}

	switch e.GetType() {
	case TypePlayerReady:
		allReady := false
		if len(g.Players) >= 5 {
			allReady = true
			for _, p := range g.Players {
				if !p.Ready {
					allReady = false
				}
			}
		}
		if allReady {
			ret = append(ret, NewGameStartEvent(len(g.Players)))
		}
	case TypePlayerAcknowledge:
		allAck := true
		for _, p := range g.Players {
			if !p.Ack {
				allAck = false
			}
		}
		if allAck {
			ret = append(ret, g.createNextRound())
		}
	case TypePlayerNominate:
		//Trigger a round.vote_start
		ret = append(ret, RoundStateEvent{
			BaseEvent: BaseEvent{Type: TypeRoundVoteStart},
			RoundID:   g.Round.ID,
			State:     RoundStateVoting,
		})
	case TypePlayerVote:
		//If all the votes are in...
		votesIn := make(map[string]bool)
		c := 0
		for _, v := range g.Round.Votes {
			votesIn[v.PlayerID] = true
			if v.Vote {
				c++
			}
		}
		allIn := true
		for _, p := range g.Players {
			if p.ExecutedBy == "" {
				if !votesIn[p.ID] {
					allIn = false
					break
				}
			}
		}
		if allIn {
			//TODO If secret hitler is elected chancellor with 3 facist polices down, facists win
			ret = append(ret, RoundVoteEndEvent{
				BaseEvent: BaseEvent{Type: TypeRoundVoteEnd},
				RoundID:   g.Round.ID,
				Votes:     g.Round.Votes,
				Succeeded: ((float64(c) / float64(len(g.Round.Votes))) * 100) > 50.0,
			})
		}
	case TypePlayerLegislate:
		//TODO Trigger a legistlate chancellor with the remaining cards or
		//TODO Trigger a legislate enact with the remaining card
	case TypePlayerInvestigate:
		//TODO Trigger a executive action enact with the revealed information
	case TypePlayerSpecialElection:
		//TODO Trigger a executive action enact that adds the special election state
	case TypePlayerExecute:
	case TypeRoundStart:
		ret = append(ret, RoundStateEvent{
			BaseEvent: BaseEvent{Type: TypeRoundNominateRequest},
			RoundID:   g.Round.ID,
			State:     RoundStateNominating,
		})
	case TypeRoundVoteEnd:
		//TODO If the vote failed, enact a policy if failed votes = 3
		if g.Board.FailedVotes >= 3 {
			ret = append(ret, RoundLegislateEnact{
				BaseEvent: BaseEvent{Type: TypeRoundLegislateEnact},
				RoundID:   g.Round.ID,
				Policy:    g.Draw[len(g.Draw)-1],
			})
		}
		//TODO Pop the top policy off the draw pile and enact it
		//TODO Should I just trigger a policy enact event?, let the logic occour in there?
	case TypeRoundLegislateEnact:
		//TODO Send an event to play the last card off the draw pile
		//TODO If the draw pile is 2 or less cards, re-shuffle
		//TODO If there was an executive action under the played policy, trigger that
	case TypeRoundExecutiveActionEnact:
		//If secret hitler is executed liberals win
		for _, p := range g.Players {
			if p.Role == RoleHitler && p.ExecutedBy != "" {
				//End the game now
				ret = append(ret, GameEndEvent{
					BaseEvent:    BaseEvent{Type: TypeGameEnd},
					WinningParty: PartyLiberal,
				})
			}
		}
	case TypeRoundEnd:
		ret = append(ret, g.createNextRound())
	}
	return ret, nil
}
