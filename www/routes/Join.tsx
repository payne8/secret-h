import * as React from 'react';
import { joinPlayer } from '../api';
import { Header } from '../components/Header';
import { appState } from '../state';

export class Join extends React.Component {
  name: string = '';

  join = () => {
    const player = {
      id: this.name.toLowerCase().replace(/\s+/, ''),
      name: this.name,
      ready: false
    };

    return joinPlayer(player.id, player.name)
      .then(() => {
        appState.setCurrentPlayer(player);
      })
      .catch(console.error);
  }

  onChange = (event: any) => {
    this.name = event.target.value;
  }

  render() {
    return (
      <div>
        <Header title="Join Game" />
        <div style={{ marginBottom: '1em' }}>
          <label>
            Name <br />
            <input type="text" name="firstName" onChange={this.onChange} />
          </label>
        </div>
        <button onClick={this.join}>
          Join
        </button>
      </div>
    );
  }
}
