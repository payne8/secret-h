import * as React from 'react';
import { render } from 'react-dom';
import { BrowserRouter, Route } from 'react-router-dom';

import { Provider, Subscribe } from 'unstated';
import { AppState } from './state';

import { Lobby } from './routes/Lobby';
import { Join } from './routes/Join';
import { MainGame } from './routes/MainGame';

class App extends React.Component<{ appState: any }> {
  componentDidMount() {
    this.props.appState
      .init() // start listenting to SSE
      .fetchInitialState();
  }

  // this is for HMR
  componentWillUnmount() {
    this.props.appState.destroy();
  }

  render() {
    return (
      <div style={{ maxWidth: 800, margin: '0 auto' }}>
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
    <Provider>
      <Subscribe to={[AppState]}>
      {(appState) => (
        <BrowserRouter>
          <App appState={appState} />
        </BrowserRouter>
      )}
      </Subscribe>
    </Provider>
  );
}


// mount it to the DOM
render(React.createElement(AppContainer), document.querySelector('#app'));
