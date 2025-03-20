// Package core contains the main business logic,
// including the domain models which implement the rule of the game and its score calculation logic.
package core

import (
	"errors"
	"fmt"
)

type Game interface {
	StartGame(playerNames []string) error
	GetCurrentFrame() int
	GetPlayers() []*Player
}

const numPin = 10
const maxPlayer = 5

type TenPinGame struct {
	players      []*Player
	currentFrame int
}

func (t *TenPinGame) GetCurrentFrame() int {
	return t.currentFrame
}

func (t *TenPinGame) StartGame(playerNames []string) error {
	if len(playerNames) == 0 {
		return errors.New("names is empty")
	}
	if len(playerNames) > maxPlayer {
		return fmt.Errorf("max num of players is %d", maxPlayer)
	}

	for i, e := range playerNames {
		if e == "" {
			return fmt.Errorf("player at index %d has empty name", i)
		}
		t.players = append(t.players, NewPlayer(e))
	}
	return nil
}

func (t *TenPinGame) GetPlayers() []*Player {
	return t.players
}

type Player struct {
	Name string
}

func NewPlayer(name string) *Player {
	return &Player{
		Name: name,
	}
}

type FrameScore interface {
	SetScores(isStrike, isSpare bool, scores ...int) error
	GetRollScores() []int
}
