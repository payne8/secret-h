let params = new URLSearchParams(location.search.slice(1));
let gameId = params.get("gameId")
let reveal = false
let gameState = {state:""}
let playerId = ""
let players = new Array()

function getStatePlayer(playerId){
	for(let p of gameState.players){
		if(p.id == playerId){
			return p
		}
	}
}

function getCachedPlayer(playerId){
	for(let p of players){
		if(p.id == playerId){
			return p
		}
	}
	return {
		id: playerId,
		email: "",
		username: playerId,
		name: playerId,
		thumbnailUri: "http://www.gravatar.com/avatar"
	}
}

getMe().then(function(me){
	if(me.err){
		//There was an error, show login/register
		console.log("me")
		console.log(me)
	}else{
		console.log("me")
		playerId = me.id
		players.push(me)
		document.querySelector("#me-banner").innerHTML = "Welcome " + me.name
		document.querySelector("#me-banner").classList.remove("no-display")
		document.querySelector("#login").classList.add("no-display")
		document.querySelector("#join").classList.remove("no-display")
	}
	initializeSSE()
}).catch(function(e){
	console.log(e)
})

function revealInfo(){
	document.querySelectorAll(".player").forEach(function(item){
		if(reveal){
			item.classList.remove("reveal-info")
		}else{
			item.classList.add("reveal-info")
		}
	})
	reveal = !reveal
}

function drawState(state){
	document.querySelector("#join").classList.add("no-display")
	document.querySelector("#ready").classList.add("no-display")
	document.querySelector("#acknowledge").classList.add("no-display")
	//TODO Only show buttons if there is an authenticated user, and only if the game is in the right state
	if(playerId != ""){
		if(state.state == ""){
			document.querySelector("#join").classList.remove("no-display")
			document.querySelector("#ready").classList.remove("no-display")
		}
		if(state.state == "init" && !getStatePlayer(playerId).ack){
			document.querySelector("#acknowledge").classList.remove("no-display")
		}
	}

	//Fill in the draw and discard piles
	while(document.querySelector("#draw").firstChild) {document.querySelector("#draw").removeChild(document.querySelector("#draw").firstChild);}
	while(document.querySelector("#discard").firstChild) {document.querySelector("#discard").removeChild(document.querySelector("#discard").firstChild);}
	if(state.draw){
		for (let p of state.draw){
			var a = document.createElement("article")
			a.innerHTML = p
			document.querySelector("#draw").appendChild(a)
		}
	}
	if(state.discard){
		for (let p of state.discard){
			var a = document.createElement("article")
			a.innerHTML = p
			document.querySelector("#discard").appendChild(a)
		}
	}

	//Fill in the liberal
	if(state.liberal > 0){ document.querySelector("#liberal1").classList.add("board-liberal-played") }
	if(state.liberal > 1){ document.querySelector("#liberal2").classList.add("board-liberal-played") }
	if(state.liberal > 2){ document.querySelector("#liberal3").classList.add("board-liberal-played") }
	if(state.liberal > 3){ document.querySelector("#liberal4").classList.add("board-liberal-played") }
	if(state.liberal > 4){ document.querySelector("#liberal5").classList.add("board-liberal-played") }
	//Fill in the failed-votes
	document.querySelector("#et0").classList.remove("board-election-tracker-spot")
	document.querySelector("#et1").classList.remove("board-election-tracker-spot")
	document.querySelector("#et2").classList.remove("board-election-tracker-spot")
	document.querySelector("#et3").classList.remove("board-election-tracker-spot")
	if(state.electionTracker == 0){ document.querySelector("#et0").classList.add("board-election-tracker-spot") }
	if(state.electionTracker == 1){ document.querySelector("#et1").classList.add("board-election-tracker-spot") }
	if(state.electionTracker == 2){ document.querySelector("#et2").classList.add("board-election-tracker-spot") }
	if(state.electionTracker == 3){ document.querySelector("#et3").classList.add("board-election-tracker-spot") }
	//Fill in the liberal
	if(state.facist> 0){ document.querySelector("#facist1").classList.add("board-facist-played") }
	if(state.facist> 1){ document.querySelector("#facist2").classList.add("board-facist-played") }
	if(state.facist> 2){ document.querySelector("#facist3").classList.add("board-facist-played") }
	if(state.facist> 4){ document.querySelector("#facist4").classList.add("board-facist-played") }
	if(state.facist> 5){ document.querySelector("#facist5").classList.add("board-facist-played") }
	if(state.facist> 6){ document.querySelector("#facist6").classList.add("board-facist-played") }

	//Add player if doesn't exist
	if(state.players){
		for (let p of state.players) {
			if(document.querySelector("#player-"+p.id) == null){
				var a = document.createElement("article")
				a.id = "player-"+p.id
				a.classList.add("player")
				var i = document.createElement("img")
				var h = document.createElement("header")
				var f = document.createElement("footer")
				i.src = getCachedPlayer(p.id).thumbnailUrl
				//Use players actual name
				h.innerHTML = getCachedPlayer(p.id).name
				a.appendChild(i)
				a.appendChild(h)
				a.appendChild(f)
				document.querySelector("#players").appendChild(a)
			}
			document.querySelector("#player-"+p.id+">footer").innerHTML = ""
			document.querySelector("#player-"+p.id+">img").src = getCachedPlayer(p.id).thumbnailUrl
			document.querySelector("#player-"+p.id+">header").innerHTML = getCachedPlayer(p.id).name

			if(state.state == ""){
				if(p.ready){ document.querySelector("#player-"+p.id+">footer").innerHTML = "Ready" }
			}
			if(state.state == "init"){
				if(p.ack){ document.querySelector("#player-"+p.id+">footer").innerHTML = "Acknowledge" }
			}
			if(state.round && state.round.presidentId == p.id){
				document.querySelector("#player-"+p.id+">footer").innerHTML = "President"
			}
			if(state.round && state.round.chancellorId == p.id){
				document.querySelector("#player-"+p.id+">footer").innerHTML = "Chancellor"
			}

			if(p.party == "facist"){document.querySelector("#player-"+p.id).classList.add("player-party-facist")}
			if(p.party == "liberal"){document.querySelector("#player-"+p.id).classList.add("player-party-liberal")}
			if(p.role== "hitler"){document.querySelector("#player-"+p.id).classList.add("player-role-hitler")}
		}
	}
	if(state.round && state.round.state){
		if(state.round.state == "voting" && state.round.votes){
			for(let v of state.round.votes){
				document.querySelector("#player-"+v.playerId+">footer").innerHTML = "Voted"
			}
		}
	}
	if(state.winningParty == "liberal"){
		document.querySelector("#win").innerHTML = "Liberals Win!"
	}
	if(state.winningParty == "facist"){
		document.querySelector("#win").innerHTML = "Facists Win!"
	}
}

function addLog(e){
}
/*
document.addEventListener('DOMContentLoaded', function() {
	drawState(s1)
}, false);
*/
