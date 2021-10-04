import React from "react";

interface EndTurnFormProps {
  gameId: string
  disabled: boolean
}

export default function EndTurnForm(props: EndTurnFormProps) {
  const {gameId, disabled} = props;

  const submit = () => {
    const event = {type: "EndTurn"};
    const body = JSON.stringify(event);
    fetch(`/api/games/${gameId}`, {method: "PUT", body})
      .then(() => console.log("Turn ended."));
  };

  return (
    <input type="button"
           id="end-turn"
           name="end-turn"
           value="End Turn"
           disabled={disabled}
           onClick={submit} />
  );
}
