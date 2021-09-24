import React, {useState} from "react";

interface OverFormProps {
  disabled: boolean
}

export default function OverForm(props: OverFormProps) {
  const {disabled} = props;

  const [over, setOver] = useState<boolean | null>(null);

  const submit = () => {
    const event = {type: "ChooseOver", over};
    const body = JSON.stringify(event);
    fetch("/api/game", {method: "PUT", body})
      .then(() => console.log("Over chosen."));
  };

  const changeTo = (over: boolean) => () => setOver(over)

  const radioStyle = {display: "inline", margin: "0.5em"};

  const submitDisabled = disabled || over === null;

  return (
    <>
      <label>
        Under / Over
      </label>
      <input type="radio"
             id="under"
             name="under"
             checked={over === false}
             onChange={changeTo(false)}
             disabled={disabled}
             style={radioStyle} />
      <label htmlFor="html"
             style={radioStyle}>
        Under
      </label>
      <input type="radio"
             id="over"
             name="over"
             checked={over === true}
             onChange={changeTo(true)}
             disabled={disabled}
             style={radioStyle} />
      <label htmlFor="html"
             style={radioStyle}>
        Over
      </label>
      <input type="button"
             id="over-submit"
             name="over-submit"
             value="Submit"
             onClick={submit}
             disabled={submitDisabled}
             style={{marginLeft: "1em"}} />
    </>
  );
}
