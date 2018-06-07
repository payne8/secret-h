function initializeSSE(){
var source = new EventSource("/api/games/"+params.get("gameId")+"/events", {
	withCredentials: true
});
source.addEventListener("state", function(e) {
	let d = JSON.parse(e.data)
	drawState(d)
	gameState = d
}, false);

source.addEventListener("player", function(e){
	let d = JSON.parse(e.data)
	players.push(d)

}, false)

source.addEventListener("player.join", function(e){
	let d = JSON.parse(e.data)
	if(d.player.id != playerId){
		getPlayer(d.player.id).then(function(p){
			players.push(p)
			//Update the player
			document.querySelector("#player-"+p.id+">img").src = p.thumbnailUrl
			document.querySelector("#player-"+p.id+">header").innerHTML = p.name
			
		})
	}
}, false)

source.addEventListener("player.ready", function(e){},false)
source.addEventListener("player.acknowledge", function(e){
	let d = JSON.parse(e.data)
	if(d.playerId == playerId){
		revealInfo()
		document.querySelector("#acknowledge").classList.add("no-display")
	}
},false)

source.addEventListener("player.nominate", function(e){
	let d = JSON.parse(e.data)
	log = document.createElement("article")
	log.innerHTML = "Player " + getCachedPlayer(d.playerId).name + " nominated " + getCachedPlayer(d.otherPlayerId).name
	document.querySelector("#log").appendChild(log)
	document.querySelector("#log").scrollTop = document.querySelector("#log").scrollHeight;
	if(d.playerId == playerId){
		console.log("removing all .request-nominate")
		document.querySelectorAll(".request-nominate").forEach(function(item){document.querySelector("#actions").removeChild(item)})
	}
}, false)

source.addEventListener("player.vote", function(e){
	let d = JSON.parse(e.data)

	document.querySelector("#player-"+d.playerId+">footer").innerHTML = "Voted"
	if(d.playerId == playerId){
		console.log("removing all .request-vote")
		document.querySelectorAll(".request-vote").forEach(function(item){document.querySelector("#actions").removeChild(item)})
	}
}, false)

source.addEventListener("player.legislate", function(e){
	let d = JSON.parse(e.data)
	log = document.createElement("article")
	if(!d.veto){
		log.innerHTML = "Player " + getCachedPlayer(d.playerId).name + " has legislated"
	}else{
		log.innerHTML = "Player " + getCachedPlayer(d.playerId).name + " has vetoed"
	}
	document.querySelector("#log").appendChild(log)
	document.querySelector("#log").scrollTop = document.querySelector("#log").scrollHeight;
	if(d.playerId == playerId){
		console.log("removing all .request-legislate")
		document.querySelectorAll(".request-legislate").forEach(function(item){document.querySelector("#actions").removeChild(item)})
	}
}, false)

source.addEventListener("player.investigate", function(e){
	let d = JSON.parse(e.data)
	log = document.createElement("article")
	log.innerHTML = "Player " + getCachedPlayer(d.playerId).name + " has investigated " + getCachedPlayer(d.otherPlayerId).name
	document.querySelector("#log").appendChild(log)
	document.querySelector("#log").scrollTop = document.querySelector("#log").scrollHeight;
	if(d.playerId == playerId){
		console.log("removing all .request-executive-action")
		document.querySelectorAll(".request-executive-action").forEach(function(item){document.querySelector("#actions").removeChild(item)})
	}
}, false)

source.addEventListener("player.special_election", function(e){
	let d = JSON.parse(e.data)
	log = document.createElement("article")
	log.innerHTML = "Player " + d.playerId + " has appointed " + d.otherPlayerId + " for the special election"
	document.querySelector("#log").appendChild(log)
	document.querySelector("#log").scrollTop = document.querySelector("#log").scrollHeight;
	if(d.playerId == playerId){
		console.log("removing all .request-executive-action")
		document.querySelectorAll(".request-executive-action").forEach(function(item){document.querySelector("#actions").removeChild(item)})
	}
}, false)

source.addEventListener("player.execute", function(e){
	let d = JSON.parse(e.data)
	log = document.createElement("article")
	log.innerHTML = "Player " + d.playerId + " has executed " + d.otherPlayerId
	document.querySelector("#log").appendChild(log)
	document.querySelector("#log").scrollTop = document.querySelector("#log").scrollHeight;
	if(d.playerId == playerId){
		console.log("removing all .request-executive-action")
		document.querySelectorAll(".request-executive-action").forEach(function(item){document.querySelector("#actions").removeChild(item)})
	}
}, false)

source.addEventListener("player.message", function(e){},false)

source.addEventListener("assert.policies", function(e){
	let d = JSON.parse(e.data)
	log = document.createElement("article")
	p = document.createElement("p")
	ul = document.createElement("ul")
	p.innerHTML = "Player " + getCachedPlayer(d.playerId).name + " asserts policies from " + d.policySource + ":"
	for(let policy of d.policies){
		li = document.createElement("li")
		li.innerHTML = policy
		ul.appendChild(li)
	}
	log.appendChild(p)
	log.appendChild(ul)
	document.querySelector("#log").appendChild(log)
	document.querySelector("#log").scrollTop = document.querySelector("#log").scrollHeight;
	if(d.playerId == playerId && d.policySource == "request.legislate"){
		console.log("removing all .request-legislate-assert")
		document.querySelectorAll(".request-legislate-assert").forEach(function(item){document.querySelector("#actions").removeChild(item)})
	}
	if(d.playerId == playerId && d.policySource == "peek"){
		console.log("removing all .game-information-assert")
		document.querySelectorAll(".game-information-assert").forEach(function(item){document.querySelector("#actions").removeChild(item)})
	}
}, false)

source.addEventListener("assert.party", function(e){
	let d = JSON.parse(e.data)
	log = document.createElement("article")
	log.innerHTML = "Player " + getCachedPlayer(d.playerId).name + " claims " + getCachedPlayer(d.otherPlayerId).name + " party is " + d.party
	document.querySelector("#log").appendChild(log)
	document.querySelector("#log").scrollTop = document.querySelector("#log").scrollHeight;
	if(d.playerId == playerId){
		console.log("removing all .game-information-assert")
		document.querySelectorAll(".game-information-assert").forEach(function(item){document.querySelector("#actions").removeChild(item)})
	}
},false)

source.addEventListener("react.player", function(e){},false)
source.addEventListener("react.event_id", function(e){},false)
source.addEventListener("react.status", function(e){},false)

source.addEventListener("request.acknowledge", function(e){
	if(playerId != ""){
		revealInfo()
		document.querySelector("#acknowledge").classList.remove("no-display")
	}
},false)

source.addEventListener("request.vote", function(e){
	let d = JSON.parse(e.data)
	console.log(d)
	if(playerId != "" && (d.playerId == "all" || d.playerId == playerId)){
		//Ask the player to vote
		a = document.createElement("article")
		a.classList.add("request-vote")
		p = document.createElement("p")
		f = function(e){
			let vote = false
			if(this.innerHTML == "yes"){
				vote=true
			}
			sendEvent(gameId, {type:"player.vote","playerId":playerId,"vote":vote}).then(function(r){
				if(r.err){
					ep = document.createElement("p")
					ep.innerHTML = r.err
					a.appendChild(ep)
				}
			})
		}
		yes = document.createElement("button")
		yes.innerHTML = "yes"
		yes.addEventListener("click", f)
		no = document.createElement("button")
		no.innerHTML = "no"
		no.addEventListener("click", f)
		p.innerHTML = "Vote yes or no for president " + getCachedPlayer(d.presidentId).name + " and chancellor " + getCachedPlayer(d.chancellorId).name
		a.appendChild(p)
		a.appendChild(yes)
		a.appendChild(no)
		document.querySelector("#actions").appendChild(a)
	}
}, false)

source.addEventListener("request.nominate", function(e){
	let d = JSON.parse(e.data)
	if(d.playerId == playerId){
		//Show a dialog asking the user to select a player as chancellor
		a = document.createElement("article")
		a.classList.add("request-nominate")
		p = document.createElement("p")
		p.innerHTML = "Select another player as chancellor"
		a.appendChild(p)
		f = function(e){
			sendEvent(gameId, {"type":"player.nominate","playerId":playerId,"otherPlayerId":this.name}).then(function(r){
				if(r.err){
					ep = document.createElement("p")
					ep.innerHTML = r.err
					a.appendChild(ep)
				}
			})
		}
		for(let p of gameState.players){
			if(p.id != playerId){
				b = document.createElement("button")
				b.innerHTML = getCachedPlayer(p.id).name
				b.name = p.id
				b.addEventListener("click", f)
				a.appendChild(b)
			}
		}
		document.querySelector("#actions").appendChild(a)
	}
}, false)


source.addEventListener("request.legislate", function(e){
	let d = JSON.parse(e.data)
	if(playerId != "" && d.playerId == playerId){
		//Ask the player to legislate
		a = document.createElement("article")
		a.classList.add("request-legislate")
		p = document.createElement("p")
		p.innerHTML = "Select a policy to discard"
		a.appendChild(p)
		f = function(e){
			sendEvent(gameId, {type:"player.legislate","playerId":playerId,"discard":this.name}).then(function(r){
				if(r.err){
					ep = document.createElement("p")
					ep.innerHTML = r.err
					a.appendChild(ep)
				}else{
					//Ask the player to assert
					a2 = document.createElement("article")
					a2.classList.add("request-legislate-assert")
					p2 = document.createElement("p")
					p2.innerHTML = "Tell other players how many liberal policies you were dealt"
					a2.appendChild(p2)
					f2 = function(e){
						apolicies = ["facist","facist","facist"]
						if(this.name == "1"){
							apolicies = ["liberal","facist","facist"]
						}else if(this.name == "2"){
							apolicies = ["liberal","liberal","facist"]
						}else if(this.name == "3"){
							apolicies = ["liberal","liberal","liberal"]
						}
						if(d.policies.length < 3){
							apolicies.splice(2,1)
						}
						sendEvent(gameId, {type:"assert.policies","playerId":playerId,"roundId":d.roundId,"token":d.token,"policySource":"request.legislate","policies":apolicies}).then(function(r){
							if(r.err){
								ep = document.createElement("p")
								ep.innerHTML = r.err
								a2.appendChild(ep)
							}
						})
					}
					b0 = document.createElement("button")
					b0.name = "0"
					b0.innerHTML = "0"
					b0.addEventListener("click",f2)
					b1 = document.createElement("button")
					b1.name = "1"
					b1.innerHTML = "1"
					b1.addEventListener("click",f2)
					b2 = document.createElement("button")
					b2.name = "2"
					b2.innerHTML = "2"
					b2.addEventListener("click",f2)
					b3 = document.createElement("button")
					b3.name = "3"
					b3.innerHTML = "3"
					b3.addEventListener("click",f2)
					a2.appendChild(b0)
					a2.appendChild(b1)
					a2.appendChild(b2)
					a2.appendChild(b3)
					document.querySelector("#actions").appendChild(a2)
				}
			})
		}
		for(let p of d.policies){
			b = document.createElement("button")
			b.innerHTML = p
			b.name = p
			b.addEventListener("click", f)
			a.appendChild(b)
		}
		document.querySelector("#actions").appendChild(a)
	}
}, false)

source.addEventListener("request.executive_action", function(e){
	let d = JSON.parse(e.data)
	if(playerId != "" && d.playerId == playerId){
		//Show a dialog asking the user to select another player to x
		a = document.createElement("article")
		a.classList.add("request-executive-action")
		p = document.createElement("p")
		p.innerHTML = "Select another player for executive action " + d.executiveAction
		f = function(e){
			sendEvent(gameId, {"type":"player."+d.executiveAction,"playerId":playerId,"otherPlayerId":this.name}).then(function(r){
				if(r.err){
					ep = document.createElement("p")
					ep.innerHTML = r.err
					a.appendChild(ep)
				}
			})
		}
		for(let p of gameState.players){
			if(p.id != playerId){
				b = document.createElement("button")
				b.innerHTML = getCachedPlayer(p.id).name
				b.name = p.id
				a.appendChild(b)
			}
		}
		document.querySelector("#actions").appendChild(a)

	}
}, false)

source.addEventListener("game.vote_results", function(e){
	let d = JSON.parse(e.data)
	log = document.createElement("article")
	p = document.createElement("p")
	ul = document.createElement("ul")
	if(d.succeeded){
		p.innerHTML = "Vote succeeded. Downvoters:"
	}else{
		p.innerHTML = "Vote failed. Downvoters:"
	}
	for(let v of d.votes){
		if(!v.vote){
			li = document.createElement("li")
			li.innerHTML = getCachedPlayer(v.playerId).name
			ul.appendChild(li)
		}
	}
	log.appendChild(p)
	log.appendChild(ul)
	document.querySelector("#log").appendChild(log)
	document.querySelector("#log").scrollTop = document.querySelector("#log").scrollHeight;
}, false)

source.addEventListener("game.information", function(e){
	let d = JSON.parse(e.data)
	if(playerId != "" && d.playerId == playerId){
		if(d.policies){
			//Ask the player to assert
			a = document.createElement("article")
			a.classList.add("game-information-assert")
			p = document.createElement("p")
			p.innerHTML = "Tell other players how many liberal policies are on the top of the draw pile"
			a.appendChild(p)
			f = function(e){
				apolicies = ["facist","facist","facist"]
				if(this.name == "1"){
					apolicies = ["liberal","facist","facist"]
				}else if(this.name == "2"){
					apolicies = ["liberal","liberal","facist"]
				}else if(this.name == "3"){
					apolicies = ["liberal","liberal","liberal"]
				}
				if(d.policies.length < 3){
					apolicies = apolicies.splice(2,1)
				}
				sendEvent(gameId, {type:"assert.policies","playerId":playerId,"roundId":d.roundId,"token":d.token,"policySource":"peek","policies":apolicies}).then(function(r){
					if(r.err){
						ep = document.createElement("p")
						ep.innerHTML = r.err
						a.appendChild(ep)
					}
				})
			}
			b0 = document.createElement("button")
			b0.name = "0"
			b0.innerHTML = "0"
			b0.addEventListener("click",f)
			b1 = document.createElement("button")
			b1.name = "1"
			b1.innerHTML = "1"
			b1.addEventListener("click",f)
			b2 = document.createElement("button")
			b2.name = "2"
			b2.innerHTML = "2"
			b2.addEventListener("click",f)
			b3 = document.createElement("button")
			b3.name = "3"
			b3.innerHTML = "3"
			b3.addEventListener("click",f)
			a.appendChild(b0)
			a.appendChild(b1)
			a.appendChild(b2)
			a.appendChild(b3)
			document.querySelector("#actions").appendChild(a)
		}else{
			a = document.createElement("article")
			a.classList.add("game-information-assert")
			p = document.createElement("p")
			p.innerHTML = "Tell other players what party " + getCachedPlayer(d.otherPlayerId).name + " is"
			a.appendChild(p)
			f = function(e){
				sendEvent(gameId, {type:"assert.party","playerId":playerId,"roundId":d.roundId,"token":d.token,"otherPlayerId":d.otherPlayerId,"party":this.name}).then(function(r){
					if(r.err){
						ep = document.createElement("p")
						ep.innerHTML = r.err
						a.appendChild(ep)
					}
				})
			}
			bl = document.createElement("button")
			bl.name = "liberal"
			bl.innerHTML = "Liberal"
			bl.addEventListener("click",f)
			bf = document.createElement("button")
			bf.name = "facist"
			bf.innerHTML = "Facist"
			bf.addEventListener("click",f)
			a.appendChild(bl)
			a.appendChild(bf)
			document.querySelector("#actions").appendChild(a)
		}
	}
}, false)

source.addEventListener("game.finished", function(e){
}, false)

source.addEventListener("server.close", function(e){
	source.close();
}, false)

window.onbeforeunload = function(){
	source.close()
}
}
