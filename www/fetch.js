function getMe(){
	return fetch("/api/players/me", {
		credentials: "same-origin",
	}).then(response => response.json())
}

function getPlayer(playerId){
	return fetch("/api/players/"+playerId, {
		credentials: "same-origin",
	}).then(response => response.json())
}

function registerPlayer(player){
	return fetch("/api/players/",{
		body: JSON.stringify(player),
		credentials: "same-origin",
		headers: {"Content-Type": "application/json"},
		method: "POST"
	}).then(response => response.json())
}

function sendEvent(gameId, e){
	return fetch("/api/games/"+gameId+"/events", {
		body: JSON.stringify(e),
		credentials: "same-origin",
		headers: {"Content-Type": "application/json"},
		method: "POST"
	}).then(response => response.json())
}

function getState(gameId){
	return fetch("/api/games/"+gameId+"/state", {
		credentials: "same-origin"
	}).then(response => response.json())
}

function createGame(){
	return fetch("/api/games/", {
		body: JSON.stringify({}),
		credentials: "same-origin",
		headers: {"Content-Type": "application/json"},
		method: "POST"
	}).then(response => response.json())
}

function getGames(){
	return fetch("/api/games/", {
		credentials: 'same-origin'
	}).then(response => response.json())
}
