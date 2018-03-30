import * as React from 'react';
import { joinPlayer, playerReady } from '../api';
import { Link } from 'react-router-dom';
import { appState } from '../state';

export class Join extends React.Component<{ history: any }> {
  name: string = '';

  join = () => {
    const player = {
      id: this.name.toLowerCase().trim().replace(/\s+/, ''),
      name: this.name
    };
    joinPlayer(player.id, player.name)
      .then(() => playerReady(player.id))
      .then(() => {
        appState.setCurrentPlayer(player);
        this.props.history.push('/game');
      })
      .catch(console.error);
  }

  onChange = (event: any) => {
    this.name = event.target.value;
  }

  render() {
    return (
      <div>
        <label>
          Name <br />
          <input type="text" name="firstName" onChange={this.onChange} />
        </label>
        <button onClick={this.join}>
          Join
        </button>
      </div>
    );
  }
}
