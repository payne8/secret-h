import * as React from 'react';
import { createPlayer } from '../api';
import { Header } from '../components/Header';
import { appState } from '../state';

export class Join extends React.Component {
  form = {
    email: '',
    password: ''
  };

  join = () => {
    return createPlayer(this.form)
      .then((player) => {
        appState.setCurrentPlayer(player);
      })
      .catch(console.error);
  }

  onChange = (event: any, key: string) => {
    this.form[key] = event.target.value;
  }

  render() {
    return (
      <form onSubmit={event => { event.preventDefault(); this.join();}}>
        <Header title="Join Game" />
        <div style={{ marginBottom: '1em' }}>
          <label style={{ marginBottom: '2em' }}>
            Name <br />
            <input type="text" name="firstName" onChange={event => this.onChange(event, 'name')} />
          </label>
          <label style={{ marginBottom: '2em' }}>
            Email <br />
            <input type="email" name="email" onChange={event => this.onChange(event, 'email')} />
          </label>
          <label>
            Password <br />
            <input type="password" name="password" onChange={event => this.onChange(event, 'password')} />
          </label>
        </div>
        <button type="submit">
          Join
        </button>
      </form>
    );
  }
}
