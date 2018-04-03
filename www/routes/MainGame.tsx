import * as React from 'react';
import { Subscribe } from 'unstated';
import { AppState, State } from '../state';
import { Player } from '../components/Player';

export class MainGame extends React.Component {
  render() {
    return (
      <Subscribe to={[AppState]}>
        {({ state }) => (
          <div>
            <h3>Players</h3>
            {state.players.map(p => (
              <Player key={p.id} name={p.name} id={p.id} />
            ))}
          </div>
        )}
      </Subscribe>
    );
  }
}
