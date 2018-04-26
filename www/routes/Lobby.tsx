import * as React from 'react';
import { If } from '../components/If';
import { Link } from 'react-router-dom';
import { appState } from '../state';
import { Subscribe } from 'unstated';
import { getGames, getGame, createGame } from '../api';
import { Async } from '../components/Async';

export class Lobby extends React.Component {
  state = {
    selectedGameId: null
  };

  selectGame = (game) => {
    this.setState({ selectedGameId: game.id });
  };

  joinGame = (game) => {
    appState.setCurrentGame(game);
  };

  createGame = async () => {
    const game = await createGame();
    this.joinGame(game);
  };

  render() {
    return (
      <section className="lobby">
        <header className="logo">
          <img src={require("../assets/sh-logo.png")} />
        </header>
        <div>
          <Async
            load={getGames}
            render={(games) => (
              <div>
                {games.map(game => (
                  <div key={game.id} style={{ background: this.state.selectedGameId === game.id ? 'yellow' : '' }} onClick={event => this.selectGame(game)}>
                    <p>{game.id}</p>
                    <p>{game.players}</p>
                  </div>
                ))}
                <button className="button" onClick={() => this.joinGame(games.find(g => g.id === this.state.selectedGameId))}>Join Game</button>
                <button className="button" onClick={() => this.createGame()}>Create Game</button>
              </div>
            )}
          />
        </div>
        {/* <div>
          <Subscribe to={[AppState]}>
            {({ state }) => (
              <div>
                <If condition={!state.initted}>
                  LOADING...
                </If>
                <If condition={state.initted}>
                  <>
                    <If condition={state.state === ''}>
                      <Link className="button" to={`/join`}>
                        Join
                      </Link>
                    </If>
                    <If condition={state.state !== 'lobby'}>
                      <Link className="button secondary" to={`/observe`}>
                        Observe
                      </Link>
                    </If>
                  </>
                </If>
              </div>
            )}
          </Subscribe>
        </div> */}
      </section>
    );
  }
}
