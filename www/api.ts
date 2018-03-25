function initSSE() {
  var source = new EventSource('sse');

  source.addEventListener('state', (e: any) => {
    console.log(e.data);
  }, false);
}

function playerJoin(id: string, name: string) {
  fetch("/api/event", {
    method: 'POST',
    body: JSON.stringify({
      "type": "player.join",
      "player": {
        "id": id,
        "name": name
      }
    }),
    headers: new Headers({ 'Content-Type': 'application/json' })
  }).then(res => res.json())
    .catch(error => console.error('Error:', error))
    .then(response => console.log('Success:', response));
}
function playerReady(id: string) {
  fetch("/api/event", {
    method: 'POST',
    body: JSON.stringify({
      "type": "player.ready",
      "player": {
        "id": id,
        "ready": true
      }
    }),
    headers: new Headers({ 'Content-Type': 'application/json' })
  }).then(res => res.json())
    .catch(error => console.error('Error:', error))
    .then(response => console.log('Success:', response));
}

function playerAcknowledge(id: string, party: string, role: string) {
  fetch("/api/event", {
    method: 'POST',
    body: JSON.stringify({
      "type": "player.acknowledge",
      "player": {
        "id": id,
        "acknowledge": true,
        "party": party,
        "role": role
      }
    }),
    headers: new Headers({ 'Content-Type': 'application/json' })
  }).then(res => res.json())
    .catch(error => console.error('Error:', error))
    .then(response => console.log('Success:', response));
}

function playerNominate(id: string, otherPlayerId: string) {
  fetch("/api/event", {
    method: 'POST',
    body: JSON.stringify({
      "type": "player.nominate",
      "playerID": id,
      "otherPlayerID": otherPlayerId
    }),
    headers: new Headers({ 'Content-Type': 'application/json' })
  }).then(res => res.json())
    .catch(error => console.error('Error:', error))
    .then(response => console.log('Success:', response));
}

function playerVote(id: string, vote: string) {
  fetch("/api/event", {
    method: 'POST',
    body: JSON.stringify({
      "type": "player.vote",
      "vote": { "playerID": id, "vote": vote }
    }),
    headers: new Headers({ 'Content-Type': 'application/json' })
  }).then(res => res.json())
    .catch(error => console.error('Error:', error))
    .then(response => console.log('Success:', response));
}
