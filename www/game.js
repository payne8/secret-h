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
		thumbnailUrl: "http://www.gravatar.com/avatar"
	}
}

getMe().then(function(me){
	if(me.err){
		//There was an error, show login/register
		console.log("me")
		console.log(me.err)
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
		if(state.state == "init"){
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
	//Fill in the fascist
	if(state.fascist> 0){ document.querySelector("#fascist1").classList.add("board-fascist-played") }
	if(state.fascist> 1){ document.querySelector("#fascist2").classList.add("board-fascist-played") }
	if(state.fascist> 2){ document.querySelector("#fascist3").classList.add("board-fascist-played") }
	if(state.fascist> 3){ document.querySelector("#fascist4").classList.add("board-fascist-played") }
	if(state.fascist> 4){ document.querySelector("#fascist5").classList.add("board-fascist-played") }
	if(state.fascist> 5){ document.querySelector("#fascist6").classList.add("board-fascist-played") }
	//Fill in the round policies
	if(state.round && state.round.policies && state.round.policies.length > 0){
		for(let i = 0; i < state.round.policies.length; i++){
			document.querySelector("#rp"+(i+1)).classList.remove("round-liberal")
			document.querySelector("#rp"+(i+1)).classList.remove("round-fascist")
			document.querySelector("#rp"+(i+1)).classList.remove("round-masked")
			document.querySelector("#rp"+(i+1)).classList.remove("no-display")
			document.querySelector("#rp"+(i+1)).classList.add("round-"+state.round.policies[i])
		}
	}else{
		document.querySelector("#rp1").classList.add("no-display")
		document.querySelector("#rp2").classList.add("no-display")
		document.querySelector("#rp3").classList.add("no-display")
	}

	//Add player if doesn't exist
	var tpl = document.querySelector("#my-player")
	if(state.players){
		for (let p of state.players) {
			if(document.querySelector("#player-"+p.id) == null){
				t = tpl.content
				t.querySelector("img").src = getCachedPlayer(p.id).thumbnailUrl
				t.querySelector("header").innerHTML = getCachedPlayer(p.id).name
				a = document.importNode(t,true)
				a = document.querySelector("#players").appendChild(a)
				a.id = "player-"+p.id
				//TODO Something better, ^ seems like when I attach I can't get a ref, and id is lost
				document.querySelector("#players").lastElementChild.id = "player-"+p.id
			}
			document.querySelector("#player-"+p.id+">footer").innerHTML = ""
			document.querySelector("#player-"+p.id+">img").src = getCachedPlayer(p.id).thumbnailUrl
			document.querySelector("#player-"+p.id+">header").innerHTML = getCachedPlayer(p.id).name

			if(state.state == ""){
				if(p.id == playerId){
					document.querySelector("#join").classList.add("no-display")
				}
				if(p.ready){
					document.querySelector("#player-"+p.id+">footer").innerHTML = "Ready"
				}
				if(p.ready && p.id == playerId){
					document.querySelector("#join").classList.add("no-display")
					document.querySelector("#ready").classList.add("no-display")
				}
			}
			if(state.state == "init"){
				if(p.ack){
					document.querySelector("#player-"+p.id+">footer").innerHTML = "Acknowledge"
				}
				if(p.id == playerId && p.ack){
					document.querySelector("#acknowledge").classList.add("no-display")
				}
			}
			if(state.state == "started" && p.id == playerId){
				document.querySelector("#message").classList.remove("no-display")
			}
			if(state.round && state.round.presidentId == p.id){
				document.querySelector("#player-"+p.id+">footer").innerHTML = "President"
			}
			if(state.round && state.round.chancellorId == p.id){
				document.querySelector("#player-"+p.id+">footer").innerHTML = "Chancellor"
			}
			if(p.executedBy != ""){
				document.querySelector("#player-"+p.id+">footer").innerHTML = "Executed"
			}

			if(p.party == "fascist"){document.querySelector("#player-"+p.id).classList.add("player-party-fascist")}
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
	if(state.winningParty == "fascist"){
		document.querySelector("#win").innerHTML = "Fascists Win!"
	}
}

function addLog(e){
}
/*
document.addEventListener('DOMContentLoaded', function() {
	drawState(s1)
}, false);
*/
