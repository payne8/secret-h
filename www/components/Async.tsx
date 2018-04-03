import * as React from 'react';
import { Spinner } from './Spinner';
import { HTTPError } from '../api';

export type AsyncStatus = 'PENDING' | 'COMPLETE' | 'REJECTED' | 'ERROR';

interface Props<T = any> {
  load: () => Promise<T>
  loader?: (state: State) => any
  error?: (state: State) => any
  render: (data: T) => any
}

interface State {
  data: any, //TODO could type this
  status: AsyncStatus,
  error: any
}

export class Async extends React.PureComponent<Props, State> {
  private abort: VoidFunction = () => {};
  private mounted = false;
  state: State = {
    data: null,
    status: 'PENDING',
    error: null
  }

  private load() {
    const ret: any = this.props.load();
    if (ret.promise) {
      return ret;
    } else {
      return { promise: ret, abort: () => { } };
    }
  }

  componentDidMount() {
    this.mounted = true;
    const { abort, promise } = this.load();
    this.abort = abort;
    promise.then((data: any) => {
      this.setState({ data, status: 'COMPLETE' });
    })
    .catch((error: HTTPError | Error) => {
      if (error instanceof HTTPError) {
        // aborts aren't an error
        if(error.status !== -1) {
          console.error(error);
          this.mounted && this.setState({ error, status: 'REJECTED' });
        } else {
          this.mounted && this.setState({ error, status: 'COMPLETE' });
        }
      }
      else {
        this.mounted && this.setState({ error, status: 'REJECTED' });
      }
    });
  }

  componentWillUnmount() {
    this.mounted = false;
    this.abort();
  }

  loader() {
    if (this.props.loader) {
      return this.props.loader(this.state);
    } else {
      return <Spinner />
    }
  }

  error() {
    if (this.props.error) {
      return this.props.error(this.state);
    } else {
      return <span>An error has occurred.</span>
    }
  }

  render() {
    return (
      <>
        {this.state.status === 'PENDING' && this.loader()}
        {this.state.status === 'REJECTED' && this.error()}
        {this.state.status === 'COMPLETE' && this.props.render(this.state.data)}
      </>
    );
  }
}
