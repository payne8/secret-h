import * as React from 'react';
import { createGame } from '../api';
import { Header } from '../components/Header';
import { appState } from '../state';

export class Create extends React.Component {
  name: string = '';

  create = () => {
    return createGame();
  }

  onChange = (event: any) => {
    this.name = event.target.value;
  }

  render() {
    return (
      <div>
        <Header title="Create Game" />
        <div style={{ marginBottom: '1em' }}>
          <label>
            Name <br />
            <input type="text" name="firstName" onChange={this.onChange} />
          </label>
        </div>
        <button onClick={this.create}>
          Create
        </button>
      </div>
    );
  }
}
