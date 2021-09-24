package main

import (
	"fmt"
)

type Hub struct {
	state      chan chan<- GameState
	clients    map[*Client]bool
	event      chan interface{}
	register   chan *Client
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		state:      make(chan chan<- GameState),
		event:      make(chan interface{}),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]map[*Client]bool),
	}
}

func findGame(games *map[string]GameState, id string) (*GameState, error) {
	game, present := games[id]

	if present {
		return (game, nil)
	}

	return (nil, errors.New("game not found"))
}

func createGame(games *map[string]GameState, id string, targetGen TargetGenerator) error {
	game, present := games[id]

	if present {
		return (
	}
}

func handleEvent(games *map[string]GameState, targetGen TargetGenerator, event interface{}) error {
	switch e := event.(type) {

	case StartGameEvent:


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

func (h *Hub) run() {
	targetGen := RandomTargetGenerator{}
	teams := []string{"a", "b"}
	games := make(map[string]GameState)

	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

		case receive := <-h.state:
			receive <- game

		case event := <-h.event:
			handleEvent(&game, &targetGen, event)

			for client := range h.clients {
				select {
				case client.send <- game:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
