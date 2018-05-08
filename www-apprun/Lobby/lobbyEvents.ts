import app from 'apprun';
import { LobbyService } from './lobbyService';
import { Game } from '../types';


class LobbyEvents {
  private lobbyService: LobbyService;

  constructor(
    private app
  ) {
    this.lobbyService = new LobbyService(app);
    let currentGame = this.lobbyService.getCurrentGame();
    let currentUser = this.lobbyService.getCurrentUser();
    this.listGames();

    if (currentGame !== null) {
      this.app.run('#updateGameId', currentGame);
    }

    if (currentUser !== null) {
      this.app.run('#updatePlayerId', currentUser);
      this.app.run('#editPlayerId', false);
    }

    this.app.on('createUser', (email, password) => {
      this.lobbyService.createPlayer(email, password);
    });

    this.app.on('createGame', () => {
      this.lobbyService.createGame().then((newGame: Game) => {
        this.listGames();
      });
    });

    this.app.on('joinGame', (id, name) => {
      this.lobbyService.joinGame(id, name);
      this.app.run('#updateGameId', id);
    });

    this.app.on('playerReady', () => {
      this.lobbyService.playerReady();
      this.app.run('#Game');
    });

    if (currentUser !== null && currentGame !== null) {
      this.app.run('joinGame', currentGame, currentUser);
    }
  }

  private listGames() {
    this.lobbyService.listGames().then((gamesList) => {
      this.app.run('#updateGamesList', gamesList);
    });
  }
}

export {
  LobbyEvents
};
