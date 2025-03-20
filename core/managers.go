package core

import (
	"errors"
	"sync/atomic"
)

var id atomic.Int32

type GameType string

const (
	TenPin GameType = "TEN_PIN"
)

type GameManager struct {
	GameById map[int32]Game
}

func NewGameManager() *GameManager {
	return &GameManager{GameById: map[int32]Game{}}
}

func (m *GameManager) StartGame(t GameType, playerNames []string) error {
	var game Game
	switch t {
	case TenPin:
		game = &TenPinGame{}
	default:
		return errors.New("game type is not supported")
	}

	if err := game.StartGame(playerNames); err != nil {
		return err
	}

	curId := id.Add(1)
	m.GameById[curId] = game
	return nil
}
