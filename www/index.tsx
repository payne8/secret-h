// The entry point

import * as React from 'react';
import { render } from 'react-dom';
import { BrowserRouter, Route } from 'react-router-dom';
import { Lobby } from './routes/Lobby';
import { Join } from './routes/Join';
import { initSSE } from './api';
initSSE(); //TODO this doesn't go here

class App extends React.Component {
  render() {
    return (
      <div>
        <Route exact path="/" component={Lobby} />
        <Route path="/join" component={Join} />
        <Route path="/observe" component={() => <span>Not implemented yet</span>} />
      </div>
    );
  }
}

const SecretH = () => (
  <BrowserRouter>
    <App />
  </BrowserRouter>
);

// mount it to the DOM
render(React.createElement(SecretH), document.querySelector('#app'));
