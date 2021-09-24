import React, {useState} from "react";

interface GuessFormProps {
  value: number | null
  disabled: boolean
}

export default function GuessForm(props: GuessFormProps) {
  const {value, disabled} = props;

  const [guess, setGuess] = useState<number | "">(value ?? "");

  const change = event => {
    const rawValue = event.target.value;
    const value = rawValue === "" ? rawValue : parseFloat(rawValue);
    setGuess(value);
  };

  const submit = () => {
    const event = {type: "MakeGuess", guess};
    const body = JSON.stringify(event);
    fetch("/api/game", {method: "PUT", body})
      .then(() => console.log("Guess made."));
  };

  return (
    <>
      <label htmlFor="guess-value">
        Guess
      </label>
      <input type="number"
             id="guess-value"
             name="guess-value"
             value={guess}
             onChange={change}
             disabled={disabled} />
      <input type="button"
             id="guess-submit"
             name="guess-submit"
             value="Submit"
             onClick={submit}
             style={{marginLeft: "1em"}}
             disabled={disabled} />
    </>
  )
}
