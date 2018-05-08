// copied from events.go
export enum Events {
  GameStateLobby = "lobby",
	GameStateInit = "init",
	GameStateStarted = "started",
	GameStateFinished = "finished",

	RoundStateNominating = "nominating",
	RoundStateVoting = "voting",
	RoundStateFailed = "failed",
	RoundStateLegislating = "legislating",
	RoundStateExecutiveAction = "executive_action",
	RoundStateFinished = "finished",

	ExecutiveActionInvestigate = "investigate",
	ExecutiveActionPeek = "peek",
	ExecutiveActionSpecialElection = "special_election",
	ExecutiveActionExecute = "execute",

	TypePlayerJoin = "player.join",
	TypePlayerReady = "player.ready",
	TypePlayerAcknowledge = "player.acknowledge",
	TypePlayerNominate = "player.nominate",
	TypePlayerVote = "player.vote",
	TypePlayerLegislate = "player.legislate",
	TypePlayerInvestigate = "player.investigate",
	TypePlayerSpecialElection = "player.special_election",
	TypePlayerExecute = "player.execute",

	TypeRequestAcknowledge = "request.acknowledge",
	TypeRequestVote = "request.vote",
	TypeRequestNominate = "request.nominate",
	TypeRequestLegislate = "request.legislate",
	TypeRequestExecutiveAction = "request.executive_action",

	TypeGameInformation = "game.information",
	TypeGameUpdate = "game.update",
}

export enum Party {
  facist = 'facist',
  liberal = 'liberal'
}

export enum Role {
  facist = 'facist',
  liberal = 'liberal',
  hitler = 'hitler'
}

export interface Game {
  [key: string]: any
}
