package main

import (
	"encoding/json"
	"testing"
)

func TestUnmarshalEvent(t *testing.T) {
	e, _ := UnmarshalEvent([]byte("{\"type\":\"MakeGuess\",\"guess\":3.125}"))

	PrintJson("Unmarshalled event", e)
}

func TestUnmarshalMakeGuessEvent(t *testing.T) {
	var res MakeGuessEvent

	json.Unmarshal([]byte("{\"guess\":3.125}"), &res)

	PrintStruct("Result", res)
}
