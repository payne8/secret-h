import * as React from 'react';
import { Provider, Subscribe, Container } from 'unstated';
import { Events } from './types';
import { getInitialState } from './api';

class SSE {
  private source: EventSource;

  init() {
    let sourceURL = process.env.NODE_ENV === 'production' ? '/sse' : 'http://localhost:8080/sse';
    this.source = new EventSource(sourceURL);
    this.source.onerror = console.error;
    this.source.onmessage = console.log;
    return this;
  }

  listen<T = object>(eventName: string, fn: (data: T, event: any) => void) {
    const callback = (event: any) => {
      let data;
      try {
        data = JSON.parse(event.data)
      }
      catch(error) {
        console.error(error);
      }
      console.log(data, event);
      fn(data, event);
    };

    this.source.addEventListener(eventName, callback);

    return () => {
      this.source.removeEventListener(eventName, callback);
    };
  }

  destroy() {
    this.source.close();
  }
}

interface Player {
  id: string
  name: string
}

interface State {
  currentPlayer: null | Player,
  players: Player[],
  state: '' | 'lobby' | 'init' | 'started' | 'finished'
}

// global app state
export class AppState extends Container<State> {
  state: State = {
    currentPlayer: null,
    players: [],
    state: ''
  };
  eventSource: SSE;

  init() {
    this.eventSource = new SSE().init();

    this.eventSource.listen<{ player: Player }>(Events.TypePlayerJoin, (data, event) => {
      this.setState({
        players: [...this.state.players, data.player]
      });
    });
  }

  async fetchInitialState() {
    const initialState = await getInitialState();
    if (initialState.players && initialState.players.length) {
      this.setState({
        players: this.state.players.concat(initialState.players)
      });
    }
    if (initialState.state) {
      this.setState({ state: initialState.state });
    }
    return this.state;
  }

  setCurrentPlayer(player: Player) {
    this.setState({ currentPlayer: player });
  }

  destroy() {
    this.eventSource.destroy();
  }
}

export const appState = new AppState();

(window as any).getState = () => {
  console.log(appState);
  return appState;
}
