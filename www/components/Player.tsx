import * as React from 'react';
import { Subscribe } from 'unstated';
import { appState } from '../state';

export class Player extends React.Component<{ name: string, id: string }> {

  render() {
    return (
      <Subscribe to={[appState as any]}>
      {({state}) => (
          <div className={`player`}>
            <div>
              {this.props.name}
            </div>
          </div>
      )}
      </Subscribe>
    );
  }
}
