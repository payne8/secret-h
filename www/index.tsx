// The entry point

import * as React from 'react';
import { render } from 'react-dom';

class App extends React.Component {
  render() {
    return (
      <div>Hey</div>
    );
  }
}

// mount it to the DOM
render(React.createElement(App), document.querySelector('#app'));
