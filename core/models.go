// Package core contains the main business logic,
// including the domain models which implement the rule of the game and its score calculation logic.
package core

import (
	"errors"
	"fmt"
)

type Game interface {
	StartGame(playerNames []string) error
	NextFrame() int
	GetCurrentFrame() int
	GetPlayers() []*Player
	SetFrameScore(playerIndex int, isStrike, isSpare bool, scores ...int) error
}

const numPin = 10
const maxPlayer = 5

type TenPinGame struct {
	players      []*Player
	currentFrame int
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

func (t *TenPinGame) GetCurrentFrame() int {
	return t.currentFrame
}

func (t *TenPinGame) NextFrame() int {
	if t.currentFrame >= 9 {
		return t.currentFrame
	}

	t.currentFrame++
	return t.currentFrame
}

func (t *TenPinGame) SetFrameScore(playerIndex int, isStrike, isSpare bool, scores ...int) error {
	if playerIndex < 0 || playerIndex >= len(t.players) {
		return errors.New("invalid player index")
	}

	return t.players[playerIndex].Frames[t.currentFrame].KnockPins(isStrike, isSpare, scores...)
}

type Player struct {
	Name   string
	Frames [10]Frame
}

func NewPlayer(name string) *Player {
	var scores [10]Frame
	for i := 0; i < 9; i++ {
		scores[i] = &NormalFrame{}
	}
	scores[9] = &LastFrame{}

	return &Player{
		Name:   name,
		Frames: scores,
	}
}

func (p *Player) GetScores() []int {
	var res []int
	for i, f := range p.Frames {
		switch frame := f.(type) {
		case *NormalFrame:
			var nextRolls []int
			j := i + 1
			for j < len(p.Frames) {
				nextRolls = append(nextRolls, p.Frames[j].GetPins()...)
				j++
			}
			res = append(res, frame.GetScore(nextRolls))
		case *LastFrame:
			res = append(res, frame.GetScore())
		}
	}
	return res
}

type Frame interface {
	KnockPins(isStrike, isSpare bool, pins ...int) error
	GetPins() []int
}

type NormalFrame struct {
	rollScores []int
}

func (n *NormalFrame) KnockPins(isStrike, isSpare bool, scores ...int) error {
	if isStrike {
		n.rollScores = []int{numPin}
		return nil
	}

	if isSpare {
		if len(scores) != 1 {
			return errors.New("invalid scores input: len must be 1 for spare")
		}
		n.rollScores = []int{scores[0], numPin - scores[0]}
		return nil
	}

	if len(scores) != 2 {
		return errors.New("invalid scores input: len must be 2 if the frame is not strike or spare")
	}
	n.rollScores = scores
	return nil
}

func (n *NormalFrame) GetPins() []int {
	return n.rollScores
}

func (n *NormalFrame) GetScore(nextRolls []int) int {
	res := 0
	for _, e := range n.rollScores {
		res += e
	}

	bonusRolls := 0
	if n.isStrike() {
		bonusRolls = 2
	} else if n.isSpare() {
		bonusRolls = 1
	}
	for i := 0; i < len(nextRolls) && i < bonusRolls; i++ {
		res += nextRolls[i]
	}

	return res
}

func (n *NormalFrame) isStrike() bool {
	return len(n.rollScores) == 1 && n.rollScores[0] == numPin
}

func (n *NormalFrame) isSpare() bool {
	return len(n.rollScores) == 2 && n.rollScores[0] != numPin && n.rollScores[0]+n.rollScores[1] == numPin
}

type LastFrame struct {
	rollScores []int
}

func (l *LastFrame) KnockPins(isStrike, isSpare bool, scores ...int) error {
	if isStrike {
		if len(scores) != 2 {
			return errors.New("invalid scores input: len must be 2 for strike of last frame")
		}
		l.rollScores = []int{numPin, scores[0], scores[1]}
		return nil
	}

	if isSpare {
		if len(scores) != 2 {
			return errors.New("invalid scores input: len must be 2 for spare")
		}
		l.rollScores = []int{scores[0], numPin - scores[0], scores[1]}
		return nil
	}

	if len(scores) != 2 {
		return errors.New("invalid scores input: len must be 2 if the frame is not strike or spare")
	}
	l.rollScores = scores
	return nil
}

func (l *LastFrame) GetPins() []int {
	return l.rollScores
}

func (l *LastFrame) GetScore() int {
	res := 0
	for _, e := range l.rollScores {
		res += e
	}
	return res
}
