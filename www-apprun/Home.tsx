import app, { Component, on } from 'apprun';

export default class HomeComponent extends Component {
  state = {
    gameId: '',
    playerId: '',
    playerIdEdit: true,
    availableGames: []
  };

  view = (state) => {
    return <div>
      {state.playerId !== '' && !state.playerIdEdit ?
        <div>
          <span>Player Name: {state.playerId}</span>&nbsp;<button onclick={e  => this.run('#editPlayerId', true)}>Change</button>
        </div>
        :
        <div>
          Player Name: <input name="playerId" placeholder="Your player id here..." value={state.playerId} /> <button onclick={e => this.run('#playerNameInput', e)}>Save Player Name</button>
        </div>
      }
      {state.gameId !== '' ?
        <div>
          <h1>Current Game: {this.cleanUpGameId(state.gameId)}</h1>
          <button onclick={e => app.run('playerReady')}>Ready?</button>
        </div>
        :
        <div>
          <h1>Create a game</h1>
          <button class="btn btn-primary" onclick={e => app.run('createGame')}>Create Game</button>
        </div>
      }
      <div>
        {state.availableGames.map((availableGame) => {
          return state.gameId === availableGame.id ?
            <div className="selected">
              <span>{this.cleanUpGameId(availableGame.id)}</span>&nbsp;<button onclick={e => this.run('#updateGameId', '')}>Change</button>
            </div>
            :
            <div>
              <span>{this.cleanUpGameId(availableGame.id)}</span>&nbsp;{state.playerId !== '' && state.gameId !== availableGame.id ? <button onclick={e => app.run('joinGame', availableGame.id, state.playerId)}>Join</button> : <button disabled="true">Join</button> }
            </div>
        })}
      </div>
    </div>
  }

  cleanUpGameId = (gameId) => {
    return gameId.split('-')[0];
  }

  update = {
    '#Home': state => state,
    '#updateGamesList': (state, availableGames) => {
      return {
        ...state,
        availableGames
      };
    },
    '#playerNameInput': (state, event) => {
      this.run('#updatePlayerId', event.target.parentElement.querySelector('input').value);
      this.run('#editPlayerId', false);
    },
    '#updatePlayerId': (state, playerId) => {
      return {
        ...state,
        playerId
      };
    },
    '#updateGameId': (state, gameId) => {
      return {
        ...state,
        gameId
      };
    },
    '#editPlayerId': (state, playerIdEdit) => {
      return {
        ...state,
        playerIdEdit
      };
    }
  }
}
