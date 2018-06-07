function joinGame(){
	sendEvent(gameId, {
		type: "player.join",
		player: {
			id: playerId
		}
	}).then(function(ret){
		console.log(ret)
	})
}

function ready(){
	sendEvent(gameId, {
		type: "player.ready",
		player: {
			id: playerId,
			ready: true
		}
	}).then(function(ret){
		console.log(ret)
	})
}

function acknowledge(){
	sendEvent(gameId, {
		type: "player.acknowledge",
		player: {
			id: playerId,
			acknowledge: true,
			party: getStatePlayer(playerId).party,
			role: getStatePlayer(playerId).role
		}
	}).then(function(ret){
		console.log(ret)
	})
}

function sendMessage(){
	sendEvent(gameId, {
		type: "player.message",
		playerId: playerId,
		message: document.querySelector("#message>input").value
	}).then(function(ret){
		message: document.querySelector("#message>input").value = ""
		console.log(ret)
	})
	return false
}
