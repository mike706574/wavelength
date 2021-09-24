package main

import (
	"encoding/json"
	"fmt"
)

type Event struct {
	Type string `json:"type"`
}

type MakeGuessEvent struct {
	GameId string
	Guess  float64 `json:"guess"`
}

type ChooseOverEvent struct {
	GameId string
	Over   bool `json:"over"`
}

type EndTurnEvent struct {
	GameId string
}

func UnmarshalEvent(in []byte) (interface{}, error) {
	e := &Event{}

	err := json.Unmarshal(in, &e)

	if err != nil {
		return nil, err
	}

	if e == nil {
		return nil, nil
	}

	var event interface{}

	switch e.Type {

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
