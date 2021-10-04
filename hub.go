package main

import (
	"errors"
	"fmt"
)

type Hub struct {
	state      chan *GameStateRequest
	states     chan *AllGameStatesRequest
	clients    map[string]map[*Client]bool
	event      chan *GameEventRequest
	register   chan *Client
	unregister chan *Client
}

type AllGameStatesRequest struct {
	receive chan<- map[string]*GameState
}

type GameStateRequest struct {
	gameId  string
	receive chan<- *GameState
}

type GameEventRequest struct {
	gameId string
	event  Event
}

func newHub() *Hub {
	return &Hub{
		state:      make(chan *GameStateRequest),
		states:     make(chan *AllGameStatesRequest),
		event:      make(chan *GameEventRequest),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]map[*Client]bool),
	}
}

func findGame(games map[string]*GameState, id string) (*GameState, error) {
	game, present := games[id]

	if present {
		return game, nil
	}

	return nil, errors.New("game not found")
}

func createGame(games map[string]*GameState, id string, targetGen TargetGenerator) error {
	_, present := games[id]

	if present {
		return errors.New("game already found")
	}

	teams := []string{"a", "b"}

	game := NewGame(targetGen, teams)

	games[id] = &game

	return nil
}

func handleEvent(games map[string]*GameState, gameId string, targetGen TargetGenerator, event Event) error {
	game, _ := findGame(games, gameId)

	switch e := event.(type) {

	case StartGameEvent:
		return createGame(games, gameId, targetGen)

	case MakeGuessEvent:
		return MakeGuess(game, e.Guess)

	case ChooseOverEvent:
		return ChooseOver(game, e.Over)

	case EndTurnEvent:
		return EndTurn(targetGen, game)

	default:
		panic(fmt.Sprintf("Unknown event type: %T", event))
	}
}

func (hub *Hub) run() {
	targetGen := &RandomTargetGenerator{}
	games := make(map[string]*GameState)

	for {
		select {
		case client := <-hub.register:
			gameId := client.gameId

			_, present := hub.clients[gameId]

			if !present {
				hub.clients[gameId] = make(map[*Client]bool)
			}

			hub.clients[gameId][client] = true

		case client := <-hub.unregister:
			gameId := client.gameId
			if _, gamePresent := hub.clients[gameId]; gamePresent {
				if _, clientPresent := hub.clients[gameId][client]; clientPresent {
					delete(hub.clients[gameId], client)
					close(client.send)
				}
			}

		case request := <-hub.states:
			request.receive <- games

		case request := <-hub.state:
			gameId := request.gameId
			game, err := findGame(games, gameId)

			if err == nil {
				request.receive <- game
			} else {
				err := createGame(games, gameId, targetGen)

				if err != nil {
					panic("what")
				}

				game, err := findGame(games, request.gameId)

				if err != nil {
					panic("whoa")
				}

				request.receive <- game
			}

		case request := <-hub.event:
			gameId := request.gameId
			event := request.event

			err := handleEvent(games, gameId, targetGen, event)

			if err == nil {
				game, err := findGame(games, gameId)

				if err == nil {
					for client := range hub.clients[gameId] {
						select {
						case client.send <- *game:
						default:
							close(client.send)
							delete(hub.clients[gameId], client)
						}
					}
				}
			}
		}
	}
}
