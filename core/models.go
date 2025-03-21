// Package core contains the main business logic,
// including the domain models which implement the rule of the game and its score calculation logic.
package core

import (
	"errors"
	"fmt"
)

// Game interface is the standard interface for all bowling games.
type Game interface {
	StartGame(playerNames []string) error
	NextFrame() int
	GetCurrentFrame() int
	GetPlayers() []*Player
	// SetFrameResult set the result of a player at a specific playerIndex in the current frame of the game
	// @params pins contains the numbers of pins knocked by each roll.
	// Examples: strike: pins = [1], non-strike: pins = [3, 4], last frame spare: pins = [4,6,5]
	SetFrameResult(playerIndex int, pins ...int) error
}

const numPin = 10
const maxPlayer = 5

// TenPinGame implements the rule of 10-pin bowling
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

func (t *TenPinGame) SetFrameResult(playerIndex int, pins ...int) error {
	if playerIndex < 0 || playerIndex >= len(t.players) {
		return errors.New("invalid player index")
	}

	return t.players[playerIndex].Frames[t.currentFrame].KnockPins(pins...)
}

// Player contains the name and roll results by frame of a player in a game
type Player struct {
	Name   string
	Frames [10]Frame
}

func NewPlayer(name string) *Player {
	var frames [10]Frame
	for i := 0; i < 9; i++ {
		frames[i] = &NormalFrame{}
	}
	frames[9] = &LastFrame{}

	return &Player{
		Name:   name,
		Frames: frames,
	}
}

func (p *Player) GetFrameResults() [][]int {
	var res [][]int
	for _, frame := range p.Frames {
		res = append(res, frame.GetPins())
	}
	return res
}

// GetScores calculates the scores of all frames
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

// Frame represents the frame result of a player
type Frame interface {
	// KnockPins set the number of pins knocked in the frame
	KnockPins(pins ...int) error
	GetPins() []int
}

// NormalFrame represents frame 1-9 in 10-pin bowling game
type NormalFrame struct {
	pins []int
}

func (n *NormalFrame) KnockPins(pins ...int) error {
	// strike
	if pins[0] == numPin {
		if len(pins) > 1 {
			return errors.New("invalid input: len must be 1 for strike")
		}
		n.pins = []int{numPin}
		return nil
	}

	if len(pins) != 2 {
		return errors.New("invalid input: len must be 2 for non-strike")
	}

	if !validPins(pins[0], pins[1]) {
		return errors.New("invalid input. sum must <= 2.")
	}

	n.pins = pins
	return nil
}

func (n *NormalFrame) GetPins() []int {
	return n.pins
}

func (n *NormalFrame) GetScore(nextRolls []int) int {
	res := 0
	for _, e := range n.pins {
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
	return len(n.pins) >= 1 && n.pins[0] == numPin
}

func (n *NormalFrame) isSpare() bool {
	return len(n.pins) >= 2 && n.pins[0]+n.pins[1] == numPin
}

// LastFrame represents the last frame (10) in 10-pin bowling game
type LastFrame struct {
	pins []int
}

func (l *LastFrame) KnockPins(pins ...int) error {
	if len(pins) < 2 {
		return errors.New("invalid input: len must be at least 2 for last frame")
	}
	if pins[0] > numPin {
		return errors.New("invalid input for first roll")
	}
	// strike
	if pins[0] == numPin {
		if len(pins) < 3 {
			return errors.New("invalid input: len must be 3 for last frame strike")
		}
		if pins[1] == numPin {
			if pins[2] > numPin {
				return errors.New("invalid input for last roll")
			}
		} else if !validPins(pins[1], pins[2]) {
			return errors.New("invalid input for 2 last roll")
		}
	} else if pins[0]+pins[1] == numPin { // spare
		if len(pins) < 3 {
			return errors.New("invalid input: len must be 3 for last frame spare")
		}
		if pins[2] > numPin {
			return errors.New("invalid input for last roll")
		}
	} else {
		if len(pins) != 2 {
			return errors.New("invalid input: len must be 2 for last frame open")
		}
	}

	l.pins = pins
	return nil
}

func (l *LastFrame) GetPins() []int {
	return l.pins
}

func (l *LastFrame) GetScore() int {
	res := 0
	for _, e := range l.pins {
		res += e
	}
	return res
}

func validPins(firstRoll, secondRoll int) bool {
	return firstRoll >= 0 && secondRoll >= 0 && firstRoll+secondRoll <= numPin
}
