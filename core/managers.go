package core

import (
	"bowling-score-tracker/configs"
	"errors"
	"sync/atomic"
)

var id atomic.Int32

type GameManager struct {
	GameById map[int32]Game
}

func NewGameManager() *GameManager {
	return &GameManager{GameById: map[int32]Game{}}
}

func (m *GameManager) StartGame(t configs.GameType, playerNames []string) (gameId int32, err error) {
	var game Game
	switch t {
	case configs.TenPin:
		game = &TenPinGame{}
	default:
		return 0, errors.New("game type is not supported")
	}

	if err := game.StartGame(playerNames); err != nil {
		return 0, err
	}

	curId := id.Add(1)
	m.GameById[curId] = game
	return curId, nil
}
