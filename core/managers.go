package core

import (
	"errors"
	"sync/atomic"

	"github.com/samber/lo"

	"bowling-score-tracker/configs"
)

var id atomic.Int32

/*
GameManager handles external requests, coordinate the domain objects and the data storage layer.
*/
type GameManager struct {
	GameById map[int32]Game
}

func NewGameManager() *GameManager {
	return &GameManager{GameById: map[int32]Game{}}
}

// GameInfo is the standard object used to communicate about the state of a game.
type GameInfo struct {
	Id           int32         `json:"id"`
	CurrentFrame int           `json:"current_frame"`
	Players      []PlayerScore `json:"players"`
}

func (m *GameManager) StartGame(t configs.GameType, playerNames []string) (g GameInfo, err error) {
	var game Game
	switch t {
	case configs.TenPin:
		game = &TenPinGame{}
	default:
		return g, errors.New("game type is not supported")
	}

	if err = game.StartGame(playerNames); err != nil {
		return g, err
	}

	curId := id.Add(1)
	m.GameById[curId] = game

	return GameInfo{
		Id:           curId,
		CurrentFrame: game.GetCurrentFrame(),
		Players:      lo.Map(game.GetPlayers(), playerToPlayerScore),
	}, nil
}

func (m *GameManager) GetGame(gameId int32) (g GameInfo, err error) {
	game := m.GameById[gameId]
	if game == nil {
		return g, errors.New("invalid game id")
	}

	return GameInfo{
		Id:           gameId,
		CurrentFrame: game.GetCurrentFrame(),
		Players:      lo.Map(game.GetPlayers(), playerToPlayerScore),
	}, nil
}

type PlayerScore struct {
	Name       string  `json:"name"`
	Frames     [][]int `json:"frames"`
	Scores     []int   `json:"scores"`
	TotalScore int     `json:"total_score"`
}

// SetFrameResult set the result of a player at a specific playerIndex in the current frame of a specific game.
// @params pins contains the numbers of pins knocked by each roll.
// Examples: strike: pins = [10], non-strike: pins = [3, 4], last frame spare: pins = [4,6,5]
func (m *GameManager) SetFrameResult(gameId int32, playerIndex int, pins ...int) (g GameInfo, err error) {
	game := m.GameById[gameId]
	if game == nil {
		return g, errors.New("invalid game id")
	}

	if err = game.SetFrameResult(playerIndex, pins...); err != nil {
		return g, err
	}

	return GameInfo{
		Id:           gameId,
		CurrentFrame: game.GetCurrentFrame(),
		Players:      lo.Map(game.GetPlayers(), playerToPlayerScore),
	}, nil
}

// NextFrame increases the current frame of a game
func (m *GameManager) NextFrame(gameId int32) (g GameInfo, err error) {
	game := m.GameById[gameId]
	if game == nil {
		return g, errors.New("invalid game id")
	}

	game.NextFrame()

	return GameInfo{
		Id:           gameId,
		CurrentFrame: game.GetCurrentFrame(),
		Players:      lo.Map(game.GetPlayers(), playerToPlayerScore),
	}, nil
}

func playerToPlayerScore(p *Player, index int) PlayerScore {
	return PlayerScore{
		Name:   p.name,
		Frames: p.GetFrameResults(),
		Scores: p.GetScores(),
		TotalScore: lo.Reduce(p.GetScores(), func(agg int, item int, index int) int {
			return agg + item
		}, 0),
	}
}
