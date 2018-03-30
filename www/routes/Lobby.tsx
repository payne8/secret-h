import * as React from 'react';
import { Async } from '../components/Async';
import { Link } from 'react-router-dom';
import { appState } from '../state';

export class Lobby extends React.Component {
  render() {
    return (
      <section className="lobby">
        <header>
          Secret H logo here
        </header>
        <div>
          <Async
            load={appState.fetchInitialState.bind(appState)}
            render={({state}) => (
              <div>
                {state === 'lobby' &&
                  <Link className="button" to={`/join`}>
                    Join
                  </Link>
                }
                {state !== 'lobby' &&
                  <Link className="button" to={`/observe`}>
                    Join
                  </Link>
                }
              </div>
            )}
          />
        </div>
      </section>
    );
  }
}
