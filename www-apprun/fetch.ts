import { Events } from './types';

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

export function postEvent(eventType: Events, payload, playerID) {
  // /api/games/{gameID}/events?playerID=${playerID}
  return fetchJSON(`/api/games/${payload.gameId}/events?playerID=${playerID}`, {
    method: 'POST',
    body: JSON.stringify({
      type: eventType,
      ...payload
    })
  });
}
