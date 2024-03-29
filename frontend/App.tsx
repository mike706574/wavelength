import React, {useEffect, useState} from "react";
import * as ReactDOM from "react-dom";

import EndTurnForm from "./EndTurnForm";
import GuessForm from "./GuessForm";
import OverForm from "./OverForm";

interface TurnRecord {
  guess: number | null
  over: boolean | null
  team: string
  target: number
}

interface GameState {
  teams: string[]
  turns: TurnRecord[]
}

function getCurrentTurn(game: GameState): TurnRecord {
  const {turns} = game;
  return turns[turns.length - 1];
}

interface GameProps {
  gameId: string
  game: GameState
}

function Term(props: {value: any}) {
  return (
    <pre className="terminal">
      {JSON.stringify(props.value, null, 2)}
    </pre>
  );
}

function Game(props: GameProps) {
  const {gameId, game} = props;

  const turn = getCurrentTurn(game);

  const {guess, over} = turn;

  const guessed = guess !== null;
  const choseOver = over !== null;

  return (
    <>
      <div className="row">
        <GuessForm gameId={gameId} value={guess} disabled={guessed} />
      </div>
      <div className="row">
        <OverForm gameId={gameId} over={over} disabled={!guessed || choseOver} />
      </div>
      <div className="row">
        <EndTurnForm gameId={gameId} disabled={!choseOver} />
      </div>
      <Term value={game} />
    </>
  );
}

async function fetchGameState(gameId: string): Promise<GameState> {
  const resp = await fetch(`/api/games/${gameId}`);
  return resp.json();
}

function connectWebsocket(gameId: string, handler: (msg: string) => void): Promise<WebSocket> {
  return new Promise((resolve, reject) => {
    const host = window.location.host;

    const url = `ws://${host}/ws/games/${gameId}`;

    const ws = new WebSocket(url);

    ws.onmessage = msg => {
      console.log("Received message.", msg);
      handler(msg)
    };

    ws.onopen = () => {
      console.log("Websocket opened.");
      resolve(ws);
    };

    ws.onerror = event => {
      console.log("Websocket error.");
      reject(event);
    };
  });
}

interface AppProps {
  gameId: string
}

export default function App(props: AppProps) {
  const {gameId} = props;

  const [game, setGame] = useState<GameState | null>(null);
  const [websocket, setWebsocket] = useState<WebSocket | null>(null);

  useEffect(async () => {
    const handler = message => {
      console.log("Received message.", message);
      const {data} = message;
      const newGame = JSON.parse(data);
      setGame(newGame);
    };

    const [game, websocket] = await Promise.all([fetchGameState(gameId),
                                                 connectWebsocket(gameId, handler)]);

    setGame(game);
    setWebsocket(websocket);
  }, []);

  if(game === null) {
    return (
      <span>Loading...</span>
    );
  }

  return (
    <div className="container">
      <div className="row">
        <div className="column"
             style={{marginTop: "1em"}}>
          <h1>Wavelength</h1>
          <Game gameId={gameId} game={game} />
        </div>
      </div>
    </div>
  );
}
