import app from 'apprun';


app.on('#', _ => app.run('#Home'));

app.on('updateGamesList', ()=> {});

app.on('//', route => {
  const menus = document.querySelectorAll('.navbar-nav li');
  for (let i = 0; i < menus.length; ++i) menus[i].classList.remove('active');
  const item = document.querySelector(`[href='${route}']`);
  item && item.parentElement.classList.add('active');
});

const view = state => <div className="container">
  <nav className="navbar navbar-expand-lg navbar-light bg-light">
    <a className="navbar-brand" href="#">Secret Hitler</a>
    <button className="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarSupportedContent" aria-controls="navbarSupportedContent"
      aria-expanded="false" aria-label="Toggle navigation">
      <span className="navbar-toggler-icon"></span>
    </button>
    <div className="collapse navbar-collapse" id="navbarSupportedContent">
      <ul className="navbar-nav mr-auto">
        <li className="nav-item active">
          <a className="nav-link" href="#Home">Lobby
            <span className="sr-only">(current)</span>
          </a>
        </li>
        <li className="nav-item">
          <a className="nav-link" href="#Game">Game</a>
        </li>
      </ul>
    </div>
  </nav>
  <div className="container" id="my-app"></div>
</div>

app.start('main', {}, view, {});

import { LobbyEvents } from './Lobby/lobbyEvents';
import Home from './Home';
import Game from './Game';

const element = 'my-app';
new Home().mount(element);
new Game().mount(element);
let lobbyEvents: LobbyEvents = new LobbyEvents(app);
