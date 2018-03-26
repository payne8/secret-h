// The entry point

import * as React from 'react';
import { render } from 'react-dom';
import { isGameLobby } from './api';
import { Async } from './components/Async';

class App extends React.Component {
  render() {
    return (
      <div>
        <header>
          Secret H logo here
        </header>
        <div>
          <Async
            load={isGameLobby}
            render={(isGameLobby) => (
              <div>
                {isGameLobby &&
                  <button>Join</button>
                }
                {!isGameLobby &&
                  <button>Observe</button>
                }
              </div>
            )}
          />
        </div>
      </div>
    );
  }
}

// mount it to the DOM
render(React.createElement(App), document.querySelector('#app'));
