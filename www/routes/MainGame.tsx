import * as React from 'react';
import { Subscribe } from 'unstated';
import { appState, State } from '../state';
import { Player } from '../components/Player';
import { If } from '../components/If';
import { playerReady } from '../api';
import { Header } from '../components/Header';

export class MainGame extends React.Component {
  state = {};

  ready = () => {
    if (appState.state.currentPlayer) {
      playerReady(appState.state.currentPlayer.id).catch(console.error);
      appState.setState({ currentPlayerReady: true });
    }
  };

  render() {
    return (
      <Subscribe to={[appState as any]}>
        {({ state }) => (
          <div>
            <Header title="Round" />

            <div>
              Game is {state.state || 'not started yet'}.
            </div>

            <div style={{ textAlign: 'center', marginBottom: '1em' }}>
              <If condition={state.state === '' && !state.currentPlayerReady}>
                <button onClick={this.ready}>I'm ready to begin</button>
                <small>(When everyone says they're ready, the game begins)</small>
              </If>
            </div>

            <div className="player-container">
              {state.players.map(p => (
                <Player key={p.id} name={p.name} id={p.id} />
              ))}
            </div>
          </div>
        )}
      </Subscribe>
    );
  }
}
