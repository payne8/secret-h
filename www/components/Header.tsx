import * as React from 'react';

export function Header(props: { title?: string }) {
  return (
    <header className="header-logo">
      <div>
        <img src={require("../assets/sh-logo.png")} />
        <h2>{props.title}</h2>
      </div>
    </header>
  );
}
