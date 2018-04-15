package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	sh "github.com/murphysean/secrethitler"
	tb "github.com/nsf/termbox-go"
	"log"
	"net/http"
	"os"
)

var (
	ghost     string
	ggameID   string
	gplayerID string
)

func init() {
	flag.StringVar(&ghost, "host", "localhost:8080", "The host")
	flag.StringVar(&ggameID, "gameid", "nil", "The gameID (required)")
	flag.StringVar(&gplayerID, "playerid", "", "The playerID")
}

func main() {
	flag.Parse()
	f, err := os.OpenFile("log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "host", "http://"+ghost)
	ctx = context.WithValue(ctx, "gameID", ggameID)
	ctx = context.WithValue(ctx, "playerID", gplayerID)

	err = tb.Init()
	if err != nil {
		log.Panicln(err)
	}
	defer tb.Close()
	tb.SetOutputMode(tb.OutputNormal)
	tb.Sync()

	currIdx := 0
	states := []sh.Game{}

	var myEvent sh.Event

	ec, sc, err := openSSE(ctx)
	if err != nil {
		log.Println("openss:", err)
		return
	}
	go func() {
		for e := range ec {
			switch e.GetType() {
			case sh.TypeRequestAcknowledge:
				myEvent = e
			case sh.TypeRequestNominate:
				re := e.(sh.RequestEvent)
				if re.PlayerID == ctx.Value("playerID").(string) {
					myEvent = e
				}
			case sh.TypeRequestVote:
				myEvent = e
			case sh.TypeRequestLegislate:
				re := e.(sh.RequestEvent)
				if re.PlayerID == ctx.Value("playerID").(string) {
					myEvent = e
				}
			case sh.TypeRequestExecutiveAction:
				re := e.(sh.RequestEvent)
				if re.PlayerID == ctx.Value("playerID").(string) {
					myEvent = e
				}
			}
		}
	}()
	go func() {
		for g := range sc {
			states = append(states, g)
			currIdx = len(states) - 1
			tb.Interrupt()
		}
	}()

	for {
		switch ev := tb.PollEvent(); ev.Type {
		case tb.EventKey:
			switch ev.Key {
			case tb.KeyEsc:
				return
			case tb.KeyArrowRight:
				//Go forward an event
				if currIdx < len(states)-1 {
					currIdx++
				}
			case tb.KeyArrowLeft:
				//Go back an event
				if currIdx > 0 {
					currIdx--
				}
			case tb.KeyArrowUp:
				//Go to oldest event
				currIdx = 0
			case tb.KeyArrowDown:
				//Go to newest event
				currIdx = len(states) - 1
			case tb.KeyCtrlJ:
				//Submit join event
				sendEvent(ctx, sh.PlayerEvent{
					BaseEvent: sh.BaseEvent{Type: sh.TypePlayerJoin},
					Player: sh.Player{
						ID: ctx.Value("playerID").(string),
					}})
			case tb.KeyCtrlR:
				//Submit ready event
				sendEvent(ctx, sh.PlayerEvent{
					BaseEvent: sh.BaseEvent{Type: sh.TypePlayerReady},
					Player: sh.Player{
						ID:    ctx.Value("playerID").(string),
						Ready: true,
					}})
			case tb.KeyCtrlA:
				//Submit ack event
				if len(states) > 0 && states[len(states)-1].State == sh.GameStateInit {
					if p, err := states[len(states)-1].GetPlayerByID(ctx.Value("playerID").(string)); err == nil {
						sendEvent(ctx, sh.PlayerEvent{
							BaseEvent: sh.BaseEvent{Type: sh.TypePlayerAcknowledge},
							Player: sh.Player{
								ID:    ctx.Value("playerID").(string),
								Party: p.Party,
								Role:  p.Role,
							}})
					}
				}
			default:
				var se sh.Event
				if myEvent != nil && len(states) > 0 {
					switch ev.Ch {
					case '0':
						if myEvent.GetType() == sh.TypeRequestNominate {
							se = sh.PlayerPlayerEvent{
								BaseEvent:     sh.BaseEvent{Type: sh.TypePlayerNominate},
								PlayerID:      ctx.Value("playerID").(string),
								OtherPlayerID: states[len(states)-1].Players[0].ID,
							}
						}
						if myEvent.GetType() == sh.TypeRequestExecutiveAction {
							se = sh.PlayerPlayerEvent{
								BaseEvent:     sh.BaseEvent{Type: "player." + myEvent.(sh.RequestEvent).ExecutiveAction},
								PlayerID:      ctx.Value("playerID").(string),
								OtherPlayerID: states[len(states)-1].Players[0].ID,
							}
						}
					case '1':
						if myEvent.GetType() == sh.TypeRequestNominate {
							se = sh.PlayerPlayerEvent{
								BaseEvent:     sh.BaseEvent{Type: sh.TypePlayerNominate},
								PlayerID:      ctx.Value("playerID").(string),
								OtherPlayerID: states[len(states)-1].Players[1].ID,
							}
						}
						if myEvent.GetType() == sh.TypeRequestExecutiveAction {
							se = sh.PlayerPlayerEvent{
								BaseEvent:     sh.BaseEvent{Type: "player." + myEvent.(sh.RequestEvent).ExecutiveAction},
								PlayerID:      ctx.Value("playerID").(string),
								OtherPlayerID: states[len(states)-1].Players[1].ID,
							}
						}
					case '2':
						if myEvent.GetType() == sh.TypeRequestNominate {
							se = sh.PlayerPlayerEvent{
								BaseEvent:     sh.BaseEvent{Type: sh.TypePlayerNominate},
								PlayerID:      ctx.Value("playerID").(string),
								OtherPlayerID: states[len(states)-1].Players[2].ID,
							}
						}
						if myEvent.GetType() == sh.TypeRequestExecutiveAction {
							se = sh.PlayerPlayerEvent{
								BaseEvent:     sh.BaseEvent{Type: "player." + myEvent.(sh.RequestEvent).ExecutiveAction},
								PlayerID:      ctx.Value("playerID").(string),
								OtherPlayerID: states[len(states)-1].Players[2].ID,
							}
						}
					case '3':
						if myEvent.GetType() == sh.TypeRequestNominate {
							se = sh.PlayerPlayerEvent{
								BaseEvent:     sh.BaseEvent{Type: sh.TypePlayerNominate},
								PlayerID:      ctx.Value("playerID").(string),
								OtherPlayerID: states[len(states)-1].Players[3].ID,
							}
						}
						if myEvent.GetType() == sh.TypeRequestExecutiveAction {
							se = sh.PlayerPlayerEvent{
								BaseEvent:     sh.BaseEvent{Type: "player." + myEvent.(sh.RequestEvent).ExecutiveAction},
								PlayerID:      ctx.Value("playerID").(string),
								OtherPlayerID: states[len(states)-1].Players[3].ID,
							}
						}
					case '4':
						if myEvent.GetType() == sh.TypeRequestNominate {
							se = sh.PlayerPlayerEvent{
								BaseEvent:     sh.BaseEvent{Type: sh.TypePlayerNominate},
								PlayerID:      ctx.Value("playerID").(string),
								OtherPlayerID: states[len(states)-1].Players[4].ID,
							}
						}
						if myEvent.GetType() == sh.TypeRequestExecutiveAction {
							se = sh.PlayerPlayerEvent{
								BaseEvent:     sh.BaseEvent{Type: "player." + myEvent.(sh.RequestEvent).ExecutiveAction},
								PlayerID:      ctx.Value("playerID").(string),
								OtherPlayerID: states[len(states)-1].Players[4].ID,
							}
						}
					case '5':
						if myEvent.GetType() == sh.TypeRequestNominate {
							se = sh.PlayerPlayerEvent{
								BaseEvent:     sh.BaseEvent{Type: sh.TypePlayerNominate},
								PlayerID:      ctx.Value("playerID").(string),
								OtherPlayerID: states[len(states)-1].Players[5].ID,
							}
						}
						if myEvent.GetType() == sh.TypeRequestExecutiveAction {
							se = sh.PlayerPlayerEvent{
								BaseEvent:     sh.BaseEvent{Type: "player." + myEvent.(sh.RequestEvent).ExecutiveAction},
								PlayerID:      ctx.Value("playerID").(string),
								OtherPlayerID: states[len(states)-1].Players[5].ID,
							}
						}
					case '6':
						if myEvent.GetType() == sh.TypeRequestNominate {
							se = sh.PlayerPlayerEvent{
								BaseEvent:     sh.BaseEvent{Type: sh.TypePlayerNominate},
								PlayerID:      ctx.Value("playerID").(string),
								OtherPlayerID: states[len(states)-1].Players[6].ID,
							}
						}
						if myEvent.GetType() == sh.TypeRequestExecutiveAction {
							se = sh.PlayerPlayerEvent{
								BaseEvent:     sh.BaseEvent{Type: "player." + myEvent.(sh.RequestEvent).ExecutiveAction},
								PlayerID:      ctx.Value("playerID").(string),
								OtherPlayerID: states[len(states)-1].Players[6].ID,
							}
						}
					case '7':
						if myEvent.GetType() == sh.TypeRequestNominate {
							se = sh.PlayerPlayerEvent{
								BaseEvent:     sh.BaseEvent{Type: sh.TypePlayerNominate},
								PlayerID:      ctx.Value("playerID").(string),
								OtherPlayerID: states[len(states)-1].Players[7].ID,
							}
						}
						if myEvent.GetType() == sh.TypeRequestExecutiveAction {
							se = sh.PlayerPlayerEvent{
								BaseEvent:     sh.BaseEvent{Type: "player." + myEvent.(sh.RequestEvent).ExecutiveAction},
								PlayerID:      ctx.Value("playerID").(string),
								OtherPlayerID: states[len(states)-1].Players[7].ID,
							}
						}
					case '8':
						if myEvent.GetType() == sh.TypeRequestNominate {
							se = sh.PlayerPlayerEvent{
								BaseEvent:     sh.BaseEvent{Type: sh.TypePlayerNominate},
								PlayerID:      ctx.Value("playerID").(string),
								OtherPlayerID: states[len(states)-1].Players[8].ID,
							}
						}
						if myEvent.GetType() == sh.TypeRequestExecutiveAction {
							se = sh.PlayerPlayerEvent{
								BaseEvent:     sh.BaseEvent{Type: "player." + myEvent.(sh.RequestEvent).ExecutiveAction},
								PlayerID:      ctx.Value("playerID").(string),
								OtherPlayerID: states[len(states)-1].Players[8].ID,
							}
						}
					case '9':
						if myEvent.GetType() == sh.TypeRequestNominate {
							se = sh.PlayerPlayerEvent{
								BaseEvent:     sh.BaseEvent{Type: sh.TypePlayerNominate},
								PlayerID:      ctx.Value("playerID").(string),
								OtherPlayerID: states[len(states)-1].Players[9].ID,
							}
						}
						if myEvent.GetType() == sh.TypeRequestExecutiveAction {
							se = sh.PlayerPlayerEvent{
								BaseEvent:     sh.BaseEvent{Type: "player." + myEvent.(sh.RequestEvent).ExecutiveAction},
								PlayerID:      ctx.Value("playerID").(string),
								OtherPlayerID: states[len(states)-1].Players[9].ID,
							}
						}
					case 'y':
						if myEvent.GetType() == sh.TypeRequestVote {
							se = sh.PlayerVoteEvent{
								BaseEvent: sh.BaseEvent{Type: sh.TypePlayerVote},
								PlayerID:  ctx.Value("playerID").(string),
								Vote:      true,
							}
						}
					case 'n':
						if myEvent.GetType() == sh.TypeRequestVote {
							se = sh.PlayerVoteEvent{
								BaseEvent: sh.BaseEvent{Type: sh.TypePlayerVote},
								PlayerID:  ctx.Value("playerID").(string),
								Vote:      false,
							}
						}

					case 'l':
						se = sh.PlayerLegislateEvent{
							BaseEvent: sh.BaseEvent{Type: sh.TypePlayerLegislate},
							PlayerID:  ctx.Value("playerID").(string),
							Discard:   sh.PolicyLiberal,
						}
					case 'f':
						se = sh.PlayerLegislateEvent{
							BaseEvent: sh.BaseEvent{Type: sh.TypePlayerLegislate},
							PlayerID:  ctx.Value("playerID").(string),
							Discard:   sh.PolicyFacist,
						}
					}
					if se != nil {
						sendEvent(ctx, se)
					}
				}
			}
		case tb.EventInterrupt:
			//Got a new state?
		}
		tb.Clear(tb.ColorDefault, tb.ColorDefault)
		if len(states) > 0 {
			drawPlayers(states[currIdx])
			drawGameBoard(states[currIdx])
			drawEventPrompt(myEvent)
		}
		tb.Flush()
	}
}

func getNameForID(id string) string {
	return id
}

func drawEventPrompt(e sh.Event) {
	if e == nil {
		return
	}
	switch e.GetType() {
	case sh.TypeRequestAcknowledge:
		drawStringAt("Ctrl-A to acknowledge your party/role:", 0, 10, tb.ColorDefault, tb.ColorDefault)
	case sh.TypeRequestNominate:
		drawStringAt("Choose another player as chancellor (0-9):", 0, 10, tb.ColorDefault, tb.ColorDefault)
	case sh.TypeRequestVote:
		drawStringAt("Vote y/n on president/chancellor:", 0, 10, tb.ColorDefault, tb.ColorDefault)
	case sh.TypeRequestLegislate:
		drawStringAt("Choose a policy to discard(l/f):", 0, 10, tb.ColorDefault, tb.ColorDefault)
	case sh.TypeRequestExecutiveAction:
		eae := e.(sh.RequestEvent)
		drawStringAt("Choose another player to "+eae.ExecutiveAction+" (0-9):", 0, 10, tb.ColorDefault, tb.ColorDefault)
	}
}

func drawPlayers(g sh.Game) {
	for i, p := range g.Players {
		var fg, bg tb.Attribute
		fg = tb.ColorDefault
		bg = tb.ColorDefault
		switch g.State {
		case sh.GameStateLobby:
			if p.Ready {
				fg = tb.ColorGreen
			}
		case sh.GameStateInit:
			if p.Ack {
				tb.SetCell(0, i, 'A', tb.ColorGreen, tb.ColorDefault)
			}
			fallthrough
		case sh.GameStateStarted:
			switch p.Role {
			case sh.RoleHitler:
				fg = tb.ColorRed | tb.AttrBold
			case sh.RoleFacist:
				fg = tb.ColorRed
			case sh.RoleLiberal:
				fg = tb.ColorBlue
			}
			if p.ExecutedBy != "" {
				tb.SetCell(0, i, 'X', tb.ColorRed, tb.ColorDefault)
			}
			if p.ID == g.Round.PresidentID {
				tb.SetCell(0, i, 'P', tb.ColorGreen, tb.ColorDefault)
			}
			if p.ID == g.Round.ChancellorID {
				tb.SetCell(0, i, 'C', tb.ColorGreen, tb.ColorDefault)
			}
		}
		drawStringAt(getNameForID(p.ID), 1, i, fg, bg)
	}
}

func drawGameBoard(g sh.Game) {
	if g.Liberal > 0 {
		tb.SetCell(20, 0, '█', tb.ColorBlue, tb.ColorDefault)
	} else {
		tb.SetCell(20, 0, '░', tb.ColorBlue, tb.ColorDefault)
	}
	if g.Liberal > 1 {
		tb.SetCell(21, 0, '█', tb.ColorBlue, tb.ColorDefault)
	} else {
		tb.SetCell(21, 0, '░', tb.ColorBlue, tb.ColorDefault)
	}
	if g.Liberal > 2 {
		tb.SetCell(22, 0, '█', tb.ColorBlue, tb.ColorDefault)
	} else {
		tb.SetCell(22, 0, '░', tb.ColorBlue, tb.ColorDefault)
	}
	if g.Liberal > 3 {
		tb.SetCell(23, 0, '█', tb.ColorBlue, tb.ColorDefault)
	} else {
		tb.SetCell(23, 0, '░', tb.ColorBlue, tb.ColorDefault)
	}
	if g.Liberal > 4 {
		tb.SetCell(24, 0, '█', tb.ColorBlue, tb.ColorDefault)
	} else {
		tb.SetCell(24, 0, '░', tb.ColorBlue, tb.ColorDefault)
	}

	tb.SetCell(21, 1, '.', tb.ColorDefault, tb.ColorDefault)
	tb.SetCell(22, 1, '.', tb.ColorDefault, tb.ColorDefault)
	tb.SetCell(23, 1, '.', tb.ColorDefault, tb.ColorDefault)
	switch g.FailedVotes {
	case 1:
		tb.SetCell(21, 1, 'x', tb.ColorDefault, tb.ColorDefault)
	case 2:
		tb.SetCell(22, 1, 'x', tb.ColorDefault, tb.ColorDefault)
	case 3:
		tb.SetCell(23, 1, 'x', tb.ColorDefault, tb.ColorDefault)
	}
	if g.WinningParty != "" {
		drawStringAt("Game Over - "+g.WinningParty+" Win!", 20, 3, tb.ColorDefault, tb.ColorDefault)
	}

	if g.Facist > 0 {
		tb.SetCell(20, 2, '█', tb.ColorRed, tb.ColorDefault)
	} else {
		tb.SetCell(20, 2, '░', tb.ColorRed, tb.ColorDefault)
	}
	if g.Facist > 1 {
		tb.SetCell(21, 2, '█', tb.ColorRed, tb.ColorDefault)
	} else {
		tb.SetCell(21, 2, '░', tb.ColorRed, tb.ColorDefault)
	}
	if g.Facist > 2 {
		tb.SetCell(22, 2, '█', tb.ColorRed, tb.ColorDefault)
	} else {
		tb.SetCell(22, 2, '░', tb.ColorRed, tb.ColorDefault)
	}
	if g.Facist > 3 {
		tb.SetCell(23, 2, '█', tb.ColorRed, tb.ColorDefault)
	} else {
		tb.SetCell(23, 2, '░', tb.ColorRed, tb.ColorDefault)
	}
	if g.Facist > 4 {
		tb.SetCell(24, 2, '█', tb.ColorRed, tb.ColorDefault)
	} else {
		tb.SetCell(24, 2, '░', tb.ColorRed, tb.ColorDefault)
	}
	if g.Facist > 5 {
		tb.SetCell(25, 2, '█', tb.ColorRed, tb.ColorDefault)
	} else {
		tb.SetCell(25, 2, '░', tb.ColorRed, tb.ColorDefault)
	}

	for i, p := range g.Draw {
		char := '?'
		fg := tb.ColorDefault
		switch p {
		case sh.PolicyFacist:
			char = 'F'
			fg = tb.ColorRed
		case sh.PolicyLiberal:
			char = 'L'
			fg = tb.ColorBlue
		}
		switch {
		case i == len(g.Draw)-1:
			tb.SetCell(27, 0, char, fg, tb.ColorDefault)
		case i == len(g.Draw)-2:
			tb.SetCell(27, 1, char, fg, tb.ColorDefault)
		case i == len(g.Draw)-3:
			tb.SetCell(27, 2, char, fg, tb.ColorDefault)
		}
	}
	for i, p := range g.Discard {
		char := '?'
		fg := tb.ColorDefault
		switch p {
		case sh.PolicyFacist:
			char = 'F'
			fg = tb.ColorRed
		case sh.PolicyLiberal:
			char = 'L'
			fg = tb.ColorBlue
		}
		switch {
		case i == len(g.Discard)-1:
			tb.SetCell(18, 0, char, fg, tb.ColorDefault)
		case i == len(g.Discard)-2:
			tb.SetCell(18, 1, char, fg, tb.ColorDefault)
		case i == len(g.Discard)-3:
			tb.SetCell(18, 2, char, fg, tb.ColorDefault)
		}
	}

	for i, p := range g.Round.Policies {
		switch p {
		case sh.PolicyFacist:
			tb.SetCell(20+i, 4, 'F', tb.ColorRed, tb.ColorDefault)
		case sh.PolicyLiberal:
			tb.SetCell(20+i, 4, 'L', tb.ColorBlue, tb.ColorDefault)
		default:
			tb.SetCell(20+i, 4, '?', tb.ColorDefault, tb.ColorDefault)
		}
	}
}

func drawStringAt(s string, x, y int, fg, bg tb.Attribute) {
	for _, r := range s {
		tb.SetCell(x, y, r, fg, bg)
		x = x + 1
	}
}

func openSSE(ctx context.Context) (<-chan sh.Event, <-chan sh.Game, error) {
	resp, err := http.Get(ctx.Value("host").(string) + "/api/games/" + ctx.Value("gameID").(string) + "/events" + "?playerID=" + ctx.Value("playerID").(string))
	if err != nil {
		return nil, nil, err
	}

	ec := make(chan sh.Event)
	gc := make(chan sh.Game)

	br := bufio.NewReader(resp.Body)
	go func() {
		defer resp.Body.Close()
		var event []byte
		for {
			b, err := br.ReadBytes('\n')
			if err != nil {
				log.Println("opensse:readbytes:", err)
				return
			}
			i := bytes.Index(b, []byte(":"))
			if i > 0 {
				switch {
				case bytes.HasPrefix(b, []byte("event: ")):
					event = bytes.TrimSpace(b[6:])
				case bytes.HasPrefix(b, []byte("data: ")):
					switch {
					case bytes.Equal(event, []byte("state")):
						g := sh.Game{}
						err := json.Unmarshal(b[5:], &g)
						if err != nil {
							continue
						}
						log.Println("opensse:sending:game:", g)
						gc <- g
					default:
						e, err := sh.UnmarshalEvent(b[5:])
						if err != nil {
							continue
						}
						log.Println("opensse:sending:event:", e)
						ec <- e
					}
				}
			}
		}
	}()

	return ec, gc, nil
}

func sendEvent(ctx context.Context, e sh.Event) error {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	err := enc.Encode(&e)
	if err != nil {
		return err
	}
	resp, err := http.Post(ctx.Value("host").(string)+"/api/games/"+ctx.Value("gameID").(string)+"/events"+"?playerID="+ctx.Value("playerID").(string), "application/json", &b)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
