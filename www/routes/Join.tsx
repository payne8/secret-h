import * as React from 'react';
import { joinPlayer } from '../api';
import { Link } from 'react-router-dom';

export class Join extends React.Component {
  name: string = '';

  join = () => {
    joinPlayer(this.name.toLowerCase().trim().replace(/\s+/, ''), this.name)
      .then(() => {
        // TODO go to MainGame or something
      })
      .catch(console.error);
  }

  onChange = (event: any) => {
    this.name = event.target.value;
  }

  render() {
    return (
      <div>
        <label>
          Name <br />
          <input type="text" name="firstName" onChange={this.onChange} />
        </label>
        <button onClick={this.join}>
          Join
        </button>
      </div>
    );
  }
}
