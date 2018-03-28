import { Events } from './types';

export function initSSE() {
  let sourceURL = process.env.NODE_ENV === 'production' ? '/sse' : 'http://localhost:8080/sse';
  const source = new EventSource(sourceURL);
  source.onerror = console.error;
  source.addEventListener('state', (e: any) => {
    console.log(e);
  }, false);
}

export class HTTPError extends Error {
  headers: object | Headers = {};
  status: number = 0;
  statusText: string = ''
  url: string = '';
  response: Response | null = null;
}

function makeHTTPError(msg: string, response: Response) {
  let err = new HTTPError(msg);
  err.headers = response.headers;
  err.status = response.status;
  err.statusText = response.statusText;
  err.url = response.url;
  err.response = response;
  return err;
}

export function fetchJSON(input: RequestInfo, init: RequestInit = {}) {
  return fetch(input, {
    ...init,
    method: init.method || 'GET',
    headers: {
      ...init.headers,
      'Content-Type': 'application/json'
    }
  })
  .then(res => {
    if (res.ok) {
      return res;
    } else {
      throw makeHTTPError(`Request to ${res.url} rejected with a status of ${res.status}`, res);
    }
  })
  .then(res => res.json().catch(err => ''));
}

export function isGameLobby(): Promise<boolean> {
  return fetchJSON('/api/state').then(data => data.state === 'lobby');
}

export function joinPlayer(id: string, name: string) {
  return fetchJSON(`api/event`, {
    method: 'POST',
    body: JSON.stringify({
      type: Events.TypePlayerJoin,
      player: {
        id,
        name
      }
    })
  });
}

export function playerReady(id: string) {
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

export function playerAcknowledge(id: string, party: string, role: string) {
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

export function playerNominate(id: string, otherPlayerId: string) {
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

export function playerVote(id: string, vote: string) {
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