import * as React from 'react';
import { render } from 'react-dom';
import { BrowserRouter, Route, withRouter } from 'react-router-dom';

import { Provider } from 'unstated';
import { appState } from './state';

import { Lobby } from './routes/Lobby';
import { Join } from './routes/Join';
import { MainGame } from './routes/MainGame';

class _App extends React.Component<{ history: any }> {
  componentDidMount() {
    appState
      .registerRouter(this.props.history);
  }

  // this is for HMR
  componentWillUnmount() {
    appState.destroy();
  }

  render() {
    return (
      <div style={{ maxWidth: 800, margin: '0 auto' }}>
        <Route exact path="/" component={Lobby} />
        <Route path="/join" component={Join} />
        <Route path="/observe" component={() => <span>Not implemented yet</span>} />
        <Route path="/game" component={MainGame} />
      </div>
    );
  }
}

const App = withRouter(_App as any);

function AppContainer() {
  return (
    <Provider inject={[appState]}>
      <BrowserRouter>
        <App />
      </BrowserRouter>
    </Provider>
  );
}


// mount it to the DOM
render(React.createElement(AppContainer), document.querySelector('#app'));
