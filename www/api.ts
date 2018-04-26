import { Events, Party, Role } from './types';
import { appState } from './state';

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
    },
    credentials: 'same-origin'
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

export function postEvent(gameId: string, playerId: string, eventType: Events, payload: object) {
  return fetchJSON(`api/games/${gameId}/events?playerId=${playerId}`, {
    method: 'POST',
    body: JSON.stringify({
      type: eventType,
      ...payload
    })
  });
}

export function getGames() {
  return fetchJSON('/api/games');
}

export function getGame(id: string) {
  return fetchJSON(`/api/games/${id}`);
}

export function createGame() {
  return fetchJSON('/api/games', { method: 'POST' });
}

export function createPlayer(player: { email: string, password: string }) {
  return fetchJSON('api/players', {
    method: 'POST',
    body: JSON.stringify(player)
  }).then(player => {
    // TODO this should be fixed server side. Server expects a number but returns a string
    player.id = parseInt(player.id);
    return player;
  });
}


// ------

async function reset() {
  // return fetch('/api/state', { method: 'PUT', body: JSON.stringify({
  //   players: []
  // }) }); //reset the game state
  return new Promise((res) => {
    fetch('/reset').catch(err => true).then(() => {
      setTimeout(res, 2000);
    })
  })
}

// async function mockGame() {
//   await reset();
//   appState.reset();
//   mockStartedGame();
// }

// (window as any).mockGame = mockGame;
// (window as any).mockStartedGame = mockStartedGame;
// (window as any).reset = reset;
// (window as any).restart = () => {
//   reset();
//   location.href = '/';
// };

// function mockStartedGame() {
//   const currentPlayer = { id: '1', name: '1', ready: false };
//   const promises = [
//     currentPlayer,
//     { id: '2', name: '2' },
//     { id: '3', name: '3' },
//     { id: '4', name: '4' },
//     { id: '5', name: '5' },
//   ].map(async (player) => {
//     let isCurrentPlayer = player.id === currentPlayer.id;

//     if (isCurrentPlayer) {
//       appState.setCurrentPlayer(currentPlayer);
//     }
//     // await joinPlayer(player.id, player.name)
//     await playerReady(player.id);
//     if (isCurrentPlayer) {
//       appState.setState({ currentPlayerReady: true });
//     }

//   });

//   return Promise.all(promises);
// }
