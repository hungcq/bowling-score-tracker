package core

import (
	"errors"
	"sync/atomic"

	"github.com/samber/lo"

	"bowling-score-tracker/configs"
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

	if err = game.StartGame(playerNames); err != nil {
		return 0, err
	}

	curId := id.Add(1)
	m.GameById[curId] = game
	return curId, nil
}

type PlayerScore struct {
	Scores     []int `json:"scores"`
	TotalScore int   `json:"total_score"`
}

func (m *GameManager) SetFrameResult(gameId int32, playerIndex int, pins ...int) ([]PlayerScore, error) {
	game := m.GameById[gameId]
	if game == nil {
		return nil, errors.New("invalid game id")
	}

	if err := game.SetFrameResult(playerIndex, pins...); err != nil {
		return nil, err
	}

	res := lo.Map(game.GetPlayers(), func(item *Player, index int) PlayerScore {
		return PlayerScore{
			Scores: item.GetScores(),
			TotalScore: lo.Reduce(item.GetScores(), func(agg int, item int, index int) int {
				return agg + item
			}, 0),
		}
	})

	return res, nil
}

func (m *GameManager) NextFrame(gameId int32) (int, error) {
	game := m.GameById[gameId]
	if game == nil {
		return 0, errors.New("invalid game id")
	}

	return game.NextFrame(), nil
}
