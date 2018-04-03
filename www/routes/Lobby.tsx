import * as React from 'react';
import { If } from '../components/If';
import { Link } from 'react-router-dom';
import { AppState, State } from '../state';
import { Subscribe } from 'unstated';

export class Lobby extends React.Component {
  render() {
    return (
      <section className="lobby">
        <header className="logo">
          <img src={require("../assets/sh-logo.png")} />
        </header>
        <div>
          <Subscribe to={[AppState]}>
            {({ state }) => (
              <div>
                <If condition={!state.initted}>
                  LOADING...
                </If>
                <If condition={state.initted}>
                  <>
                    <div style={{ marginBottom: '2em', textAlign: 'center' }}>
                      {state.players.length} players joined
                    </div>
                    <If condition={state.state === 'lobby'}>
                      <Link className="button" to={`/join`}>
                        Join
                      </Link>
                    </If>
                    <If condition={state.state !== 'lobby'}>
                      <Link className="button" to={`/observe`}>
                        Observe
                      </Link>
                    </If>
                  </>
                </If>
              </div>
            )}
          </Subscribe>
        </div>
      </section>
    );
  }
}
