import app from 'apprun';
import { fetchJSON, postEvent } from '../fetch';
import { Events } from '../types';

class LobbyService {
  private currentGameId: string;
  private playerId: string;

  constructor(
    private app
  ) {
    this.currentGameId = localStorage.getItem('currentGameId');
    this.playerId = localStorage.getItem('playerId');
  }

  public getCurrentGame() {
    return this.currentGameId;
  }

  public getCurrentUser() {
    return this.playerId;
  }

  public createPlayer(email: string, password: string) {
    return fetchJSON('/api/players', {
      method: 'POST',
      body: JSON.stringify({
          email,
          password
      })
    });
  }

  public listGames() {
    return fetchJSON('/api/games');
  }

  public createGame() {
    return fetchJSON('/api/games', {method: 'POST'}).then((res) => {
      return res;
    });
  }

  public joinGame(id: string, name: string) {
    this.currentGameId = id;
    this.playerId = name;
    localStorage.setItem('currentGameId', id);
    localStorage.setItem('playerId', name);
    return postEvent(Events.TypePlayerJoin, {
      gameId: this.currentGameId,
      playerID: this.playerId,
      player: {
        id: this.playerId
      }
    }, this.playerId).catch((errRes) => {});
  }

  public playerReady() {
    let id = this.currentGameId;
    return postEvent(Events.TypePlayerReady, {
      gameId: this.currentGameId,
      player: {
        id: this.playerId,
        ready: true
      }
    }, this.playerId);
  }
}

export {
  LobbyService
};
