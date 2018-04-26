import * as React from 'react';
import { Provider, Subscribe, Container } from 'unstated';
import { Events, Party, Role } from './types';
import { getGames, postEvent } from './api';

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
  email?: string
  username?: string
  name?: string
  thumbnailUrl?: string
  ready: boolean
  party?: Party
  role?: Role
}

export interface Game {
  id: string
  state: '' | 'lobby' | 'init' | 'started' | 'finished'
}

export interface State {
  initted: boolean
  currentPlayer: null | Player
  game: null | Game
}

// global app state
export class AppState extends Container<State> {
  state: State = {
    initted: false,
    currentPlayer: null,
    game: null
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
    });

    this.eventSource.listen<any>('state', (state) => {
      console.log(`state`, state);
    });

    return this;
  }

  reset() {
    this.setState({
      initted: false,
      currentPlayer: null,
      game: null
    });
  }

  async setCurrentPlayer(player: Player) {
    this.setState({ currentPlayer: player });
    await postEvent(this.state.game.id, player.id, Events.TypePlayerJoin, this.state.currentPlayer);
    this.router.push('/game');
  }

  async setCurrentGame(game) {
    this.setState({ game });
    if (!this.state.currentPlayer) {
      this.router.push('/join');
    } else {
      await postEvent(game.id, this.state.currentPlayer.id, Events.TypePlayerJoin, this.state.currentPlayer);
      this.router.push('/game');
    }
  }

  destroy() {
    this.eventSource.destroy();
  }
}

export const appState = new AppState();

(window as any).getState = () => appState;
