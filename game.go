package main

import (
	"errors"
	"math/rand"
)

// Types
type OverUnder int

type TurnRecord struct {
	Team   string   `json:"team"`
	Target float64  `json:"target"`
	Guess  *float64 `json:"guess"`
	Over   *bool    `json:"over"`
}

type GameState struct {
	Turns []TurnRecord `json:"turns"`
	Teams []string     `json:"teams"`
}

// Target Generator
type TargetGenerator interface {
	Target() float64
}

// Constant Target Generator
type ConstantTargetGenerator struct {
	value float64
}

func (targetGen ConstantTargetGenerator) Target() float64 {
	return targetGen.value
}

// Sequential Target Generator
type SequentialTargetGenerator struct {
	sequence []float64
	index    int
}

func (targetGen *SequentialTargetGenerator) Target() float64 {
	index := targetGen.index
	sequence := targetGen.sequence

	targetGen.index += 1

	if targetGen.index == len(sequence) {
		targetGen.index = 0
	}

	return targetGen.sequence[index]
}

// Random Target Generator
type RandomTargetGenerator struct{}

func (targetGen *RandomTargetGenerator) Target() float64 {
	return rand.Float64() * 100.0
}

// Functions
func NewGame(targetGen TargetGenerator, teams []string) GameState {
	if len(teams) < 2 {
		panic("expected at least two teams")
	}

	team := teams[0]
	target := targetGen.Target()
	turn := TurnRecord{Team: team, Target: target}

	return GameState{Teams: teams, Turns: []TurnRecord{turn}}
}

func GetCurrentTurn(game *GameState) *TurnRecord {
	turns := game.Turns
	return &turns[len(turns)-1]
}

func MakeGuess(game *GameState, guess float64) error {
	turn := GetCurrentTurn(game)

	if turn.Guess != nil {
		return errors.New("already guessed")
	}

	turn.Guess = &guess

	return nil
}

func ChooseOver(game *GameState, over bool) error {
	turn := GetCurrentTurn(game)

	if turn.Guess == nil {
		return errors.New("no guess")
	}

	turn.Over = &over

	return nil
}

func indexOf(strs []string, target string) int {
	for idx, str := range strs {
		if str == target {
			return idx
		}
	}
	return -1
}

func EndTurn(targetGen TargetGenerator, game *GameState) error {
	currentTurn := GetCurrentTurn(game)

	if currentTurn.Guess == nil {
		return errors.New("no guess")
	}

	if currentTurn.Over == nil {
		return errors.New("no guess")
	}

	currentTeam := currentTurn.Team

	teams := game.Teams

	currentTeamIndex := indexOf(teams, currentTeam)

	if currentTeamIndex == -1 {
		panic("team not found")
	}

	teamIndex := currentTeamIndex + 1

	if teamIndex == len(teams) {
		teamIndex = 0
	}

	team := teams[teamIndex]
	target := targetGen.Target()
	turn := TurnRecord{Team: team, Target: target}

	game.Turns = append(game.Turns, turn)

	return nil
}
