// The entry point

import * as React from 'react';
import { render } from 'react-dom';
import { BrowserRouter, Route } from 'react-router-dom';

import { Provider } from 'unstated';
import { appState } from './state';

import { Lobby } from './routes/Lobby';
import { Join } from './routes/Join';
import { MainGame } from './routes/MainGame';

class App extends React.Component {
  componentDidMount() {
    appState.init(); // start listenting to SSE
  }

  // this is for HMR
  componentWillUnmount() {
    appState.destroy();
  }

  render() {
    return (
      <div>
        <Route exact path="/" component={Lobby} />
        <Route path="/join" render={({history}) => <Join history={history} />} />
        <Route path="/observe" component={() => <span>Not implemented yet</span>} />
        <Route path="/game" component={MainGame} />
      </div>
    );
  }
}

function AppContainer() {
  return (
    <Provider
      inject={[appState]}
    >
      <BrowserRouter>
        <App />
      </BrowserRouter>
    </Provider>
  );
}


// mount it to the DOM
render(React.createElement(AppContainer), document.querySelector('#app'));
