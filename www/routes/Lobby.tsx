import * as React from 'react';
import { If } from '../components/If';
import { Link } from 'react-router-dom';
import { appState, AppState, State } from '../state';
import { Subscribe } from 'unstated';

export class Lobby extends React.Component {
  chooseGame(game: any) {
    appState.setCurrentGame(game);
  }

  createGame() {
    appState.createGame();
  }

  render() {
    return (
      <section className="lobby">
        <header className="logo">
          <img src={require("../assets/sh-logo.png")} />
        </header>
        <div>
          <Subscribe to={[AppState]}>
            {({ state }) => (
              <div>
                <If condition={!state.initted}>
                  LOADING...
                </If>
                <If condition={state.initted}>
                  <>
                    <If condition={state.state === ''}>
                      {state.availableGames.map((avalableGame, i) => {
                        return <Link className="button" key={avalableGame.id} to={`/join`} onClick={() => {this.chooseGame(avalableGame)}>
                          Join {avalableGame.id.split('-')[0]}
                        </Link>
                      })}
                      <button className="button" onClick={() => {this.createGame()}}>
                        Create
                      </button>
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
        </div>
      </section>
    );
  }
}
