import * as React from 'react';

export class Player extends React.Component<{ name: string, id: string }> {
  render() {
    return (
      <div>
        <div>
          {this.props.name}
        </div>
      </div>
    );
  }
}
