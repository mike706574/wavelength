import React from "react";
import * as ReactDOM from "react-dom";

import App from "./App";

function start() {
  const path = window.location.pathname;

  const gameId = path.charAt(0) === "/" ? path.substr(1) : path;

  ReactDOM.render(
    <App gameId={gameId} />,
    document.getElementById('root')
  );
}

start();
