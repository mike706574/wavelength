package main

import (
	"encoding/json"
	"fmt"
)

type Typed struct {
	Type string `json:"type"`
}

type Event interface{}

type StartGameEvent struct{}

type MakeGuessEvent struct {
	Guess float64 `json:"guess"`
}

type ChooseOverEvent struct {
	Over bool `json:"over"`
}

type EndTurnEvent struct{}

// Unmarshalling events
func UnmarshalEvent(in []byte) (Event, error) {
	e := &Typed{}

	err := json.Unmarshal(in, &e)

	if err != nil {
		return nil, err
	}

	if e == nil {
		return nil, nil
	}

	var event Event

	switch e.Type {

	case "StartGame":
		var startGameEvent StartGameEvent
		json.Unmarshal(in, &startGameEvent)
		event = startGameEvent

	case "MakeGuess":
		var makeGuessEvent MakeGuessEvent
		json.Unmarshal(in, &makeGuessEvent)
		event = makeGuessEvent

	case "ChooseOver":
		var chooseOverEvent ChooseOverEvent
		json.Unmarshal(in, &chooseOverEvent)
		event = chooseOverEvent

	case "EndTurn":
		var endTurnEvent EndTurnEvent
		json.Unmarshal(in, &endTurnEvent)
		event = endTurnEvent

	default:
		panic(fmt.Sprintf("invalid event type: %s", e.Type))
	}

	return event, nil
}
