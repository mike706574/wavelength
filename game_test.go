package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGameJson(t *testing.T) {
	guess := 73.51
	turn := TurnRecord{Team: "a", Target: 4.5, Guess: &guess}

	json, err := json.Marshal(turn)

	if err != nil {
		assert.Fail(t, "failed to marshal")
	}

	fmt.Println(string(json))
}

func TestConstantTargetGenerator(t *testing.T) {
	targetGen := ConstantTargetGenerator{5.23}

	assert.Equal(t, 5.23, targetGen.Target())
}

func TestRandomTargetGenerator(t *testing.T) {
	targetGen := RandomTargetGenerator{}

	target := targetGen.Target()

	assert.True(t, target > 0)
	assert.True(t, target < 100)
}

func TestNewGame(t *testing.T) {
	target := 42.11

	targetGen := ConstantTargetGenerator{target}

	teams := []string{"a", "b"}
	game := NewGame(&targetGen, teams)

	assert.Equal(t, teams, game.Teams)

	turns := game.Turns
	assert.Equal(t, 1, len(turns))

	turn := turns[0]

	assert.Equal(t, "a", turn.Team)
	assert.Equal(t, target, turn.Target)
	assert.Nil(t, turn.Guess)
	assert.Nil(t, turn.Over)
}

func TestGameplay(t *testing.T) {
	targetGen := SequentialTargetGenerator{sequence: []float64{42.11, 24.12}}
	teams := []string{"a", "b"}

	// New game
	game := NewGame(&targetGen, teams)

	// Make a guess
	turn := GetCurrentTurn(&game)

	MakeGuess(&game, 63.09)

	assert.Equal(t, 1, len(game.Turns))
	assert.Equal(t, 63.09, *turn.Guess)
	assert.Nil(t, turn.Over)

	// Choose over
	ChooseOver(&game, true)

	assert.Equal(t, 1, len(game.Turns))
	assert.True(t, *turn.Over)

	// End turn
	EndTurn(&targetGen, &game)

	assert.Equal(t, 2, len(game.Turns))

	turn = GetCurrentTurn(&game)

	assert.Equal(t, "b", turn.Team)
	assert.Equal(t, 24.12, turn.Target)
	assert.Nil(t, turn.Guess)
	assert.Nil(t, turn.Over)
}
