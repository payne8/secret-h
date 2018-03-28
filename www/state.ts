import * as React from 'react';
import { Provider, Subscribe, Container } from 'unstated';

type State = {
  player: {
    id: string,
    name: string
  }
};

export class AppState extends Container<State> {
  state: State = {
    player: {
      id: '',
      name: ''
    }
  };

  setPlayer(player: State['player']) {
    this.setState({ player });
  }
}

export const appState = new AppState();
