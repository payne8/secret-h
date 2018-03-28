import * as React from 'react';
import { Async } from '../components/Async';
import { isGameLobby } from '../api';
import { Link } from 'react-router-dom';

export class Lobby extends React.Component {
  render() {
    return (
      <section className="lobby">
        <header>
          Secret H logo here
        </header>
        <div>
          <Async
            load={isGameLobby}
            render={(isGameLobby) => (
              <div>
                {isGameLobby &&
                  <Link className="button" to={`/join`}>
                    Join
                  </Link>
                }
                {!isGameLobby &&
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
