import * as React from 'react';
import { Provider, Subscribe, Container } from 'unstated';
import { Events, Party, Role } from './types';
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

export interface Player {
  id: string
  name: string
  ready: boolean
  party?: Party
  role?: Role
}

export interface State {
  currentPlayer: Player
  currentPlayerReady: boolean
  players: Player[]
  state: '' | 'lobby' | 'init' | 'started' | 'finished'
  initted: boolean
}

// global app state
export class AppState extends Container<State> {
  state: State = {
    currentPlayer: { id: '1', name: 'Default', ready: false },
    currentPlayerReady: false,
    players: [],
    state: '',
    initted: false
  };
  eventSource: SSE;
  router;

  registerRouter(router) {
    this.router = router;
    return this;
  }

  init() {
    this.eventSource = new SSE().init();

    this.eventSource.listen<any>(Events.TypeGameUpdate, (state) => {
      console.log(`game update`, state);
      this.setState({ ...this.state, ...state.game });
      if(state.game.nextPresidentID === this.state.currentPlayer.id) {

      }
    });

    this.eventSource.listen<any>('state', (state) => {
      console.log(`state`, state);
    });

    return this;
  }

  reset() {
    this.setState({
      currentPlayer: { id: '1', name: 'Default', ready: false },
      currentPlayerReady: false,
      players: [],
      state: '',
      initted: false
    });
    return this.fetchInitialState();
  }

  async fetchInitialState() {
    const initialState = await getInitialState();
    this.setState({ ...this.state, ...initialState, initted: true });
    console.log(this.state);
    return this.state;
  }

  setCurrentPlayer(player: Player) {
    this.setState({ currentPlayer: player });
    this.router.push('/game');
  }

  destroy() {
    this.eventSource.destroy();
  }
}

export const appState = new AppState();

(window as any).getState = () => appState;
