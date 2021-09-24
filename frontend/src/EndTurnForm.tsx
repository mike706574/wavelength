import React from "react";

interface EndTurnFormProps {
  disabled: boolean
}

export default function EndTurnForm(props: EndTurnFormProps) {
  const {disabled} = props;

  const submit = () => {
    const event = {type: "EndTurn"};
    const body = JSON.stringify(event);
    fetch("/api/game", {method: "PUT", body})
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
