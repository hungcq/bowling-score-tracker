package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTenPinGame(t *testing.T) {
	t.Run("StartGame", func(t *testing.T) {
		t.Run("should_accept_1_player", func(t *testing.T) {
			game := &TenPinGame{}
			assert.NoError(t, game.StartGame([]string{"hung"}))
		})

		t.Run("should_accept_up_to_5_players", func(t *testing.T) {
			game := &TenPinGame{}
			assert.NoError(t, game.StartGame([]string{"hung1", "hung2", "hung3", "hung4", "hung5"}))
		})

		t.Run("should_accept_duplicated_player_names", func(t *testing.T) {
			game := &TenPinGame{}
			assert.NoError(t, game.StartGame([]string{"hung", "hung", "hung"}))
		})

		t.Run("should_accept_duplicated_player_names", func(t *testing.T) {
			game := &TenPinGame{}
			assert.NoError(t, game.StartGame([]string{"hung", "hung", "hung"}))
		})

		t.Run("should_reject_empty_names", func(t *testing.T) {
			game := &TenPinGame{}
			assert.Error(t, game.StartGame([]string{"hung", "", "hung"}))
		})

		t.Run("should_reject_empty_names_array", func(t *testing.T) {
			game := &TenPinGame{}
			assert.Error(t, game.StartGame(nil))
		})

		t.Run("should_reject_when_there_is_more_than_max_players", func(t *testing.T) {
			game := &TenPinGame{}
			var names []string
			for i := 0; i < maxPlayer+1; i++ {
				names = append(names, "hung")
			}
			assert.Error(t, game.StartGame(names))
		})
	})

	t.Run("GetPlayers", func(t *testing.T) {
		t.Run("when_no_players_added", func(t *testing.T) {
			game := &TenPinGame{}
			players := game.GetPlayers()
			assert.Empty(t, players, "expected no players initially")
		})

		t.Run("when_players_are_added", func(t *testing.T) {
			game := &TenPinGame{}
			err := game.StartGame([]string{"player1", "player2"})
			require.NoError(t, err)

			players := game.GetPlayers()

			assert.Len(t, players, 2, "expected two players after starting game")
		})
	})

	t.Run("NextFrame", func(t *testing.T) {
		t.Run("when_increase_frame_below_9", func(t *testing.T) {
			game := &TenPinGame{currentFrame: 4}

			game.NextFrame()

			assert.Equal(t, 5, game.GetCurrentFrame(), "should_increment_frame")
		})

		t.Run("when_frame_is_last", func(t *testing.T) {
			game := &TenPinGame{currentFrame: 9}

			game.NextFrame()

			assert.Equal(t, 9, game.GetCurrentFrame(), "should_not_increment_frame")
		})
	})

	t.Run("SetFrameScore", func(t *testing.T) {
		t.Run("should_reject_invalid_player_index", func(t *testing.T) {
			game := &TenPinGame{}
			err := game.StartGame([]string{"hung"})
			require.NoError(t, err)

			err = game.SetFrameScore(-1, false, false, 3, 4)
			assert.Error(t, err, "should error for invalid player index")

			err = game.SetFrameScore(1, false, false, 3, 4)
			assert.Error(t, err, "should error for invalid player index when out of bounds")
		})

		t.Run("should_reject_invalid_scores_input", func(t *testing.T) {
			t.Run("for_normal_frame", func(t *testing.T) {
				game := &TenPinGame{}
				err := game.StartGame([]string{"hung"})
				require.NoError(t, err)
				game.currentFrame = 0

				err = game.SetFrameScore(0, false, false, 7)
				assert.NotNil(t, err, "open frame require 2 scores")

				err = game.SetFrameScore(0, false, true)
				assert.NotNil(t, err, "spare frame require 1 scores")
			})

			t.Run("for_last_frame", func(t *testing.T) {
				game := &TenPinGame{}
				err := game.StartGame([]string{"hung"})
				require.NoError(t, err)
				game.currentFrame = 9

				err = game.SetFrameScore(0, false, false, 2, 2, 3)
				assert.NotNil(t, err, "open frame require 2 scores")

				err = game.SetFrameScore(0, true, false, 1)
				assert.NotNil(t, err, "strike frame require 2 scores")

				err = game.SetFrameScore(0, true, false, 1)
				assert.NotNil(t, err, "spare frame require 2 scores")
			})
		})
		t.Run("normal_frame_strike", func(t *testing.T) {
			game := &TenPinGame{}
			err := game.StartGame([]string{"hung"})
			require.NoError(t, err)
			game.currentFrame = 0

			err = game.SetFrameScore(0, true, false)

			assert.NoError(t, err, "should return success")
			rolls := game.GetPlayers()[0].Frames[game.currentFrame].GetPins()
			assert.Equal(t, []int{numPin}, rolls, "should record a single roll with all pins knocked down")
		})

		t.Run("normal_frame_spare", func(t *testing.T) {
			game := &TenPinGame{}
			err := game.StartGame([]string{"hung"})
			require.NoError(t, err)
			game.currentFrame = 1

			err = game.SetFrameScore(0, false, true, 4)

			assert.NoError(t, err, "should return success")
			rolls := game.GetPlayers()[0].Frames[game.currentFrame].GetPins()
			expected := []int{4, numPin - 4}
			assert.Equal(t, expected, rolls, "should record two rolls summing to all pins")
		})

		t.Run("normal_frame_open", func(t *testing.T) {
			game := &TenPinGame{}
			err := game.StartGame([]string{"hung"})
			require.NoError(t, err)
			game.currentFrame = 1

			err = game.SetFrameScore(0, false, false, 3, 5)

			assert.NoError(t, err, "should return success")
			rolls := game.GetPlayers()[0].Frames[game.currentFrame].GetPins()
			expected := []int{3, 5}
			assert.Equal(t, expected, rolls, "open frame should record the provided two roll scores")
		})

		t.Run("last_frame_strike", func(t *testing.T) {
			game := &TenPinGame{}
			err := game.StartGame([]string{"hung"})
			require.NoError(t, err)
			game.currentFrame = 9

			err = game.SetFrameScore(0, true, false, 7, 2)
			assert.NoError(t, err, "should return success")
			rolls := game.GetPlayers()[0].Frames[game.currentFrame].GetPins()
			expected := []int{numPin, 7, 2}
			assert.Equal(t, expected, rolls, "should record three rolls")
		})

		t.Run("last_frame_spare", func(t *testing.T) {
			game := &TenPinGame{}
			err := game.StartGame([]string{"hung"})
			require.NoError(t, err)
			game.currentFrame = 9

			err = game.SetFrameScore(0, false, true, 6, 8)
			assert.NoError(t, err, "should return success")
			rolls := game.GetPlayers()[0].Frames[game.currentFrame].GetPins()
			expected := []int{6, numPin - 6, 8}
			assert.Equal(t, expected, rolls, "should record three rolls (including bonus)")
		})

		t.Run("last_frame_open", func(t *testing.T) {
			game := &TenPinGame{}
			err := game.StartGame([]string{"hung"})
			require.NoError(t, err)
			game.currentFrame = 9

			err = game.SetFrameScore(0, false, false, 3, 6)
			assert.NoError(t, err)
			rolls := game.GetPlayers()[0].Frames[game.currentFrame].GetPins()
			expected := []int{3, 6}
			assert.Equal(t, expected, rolls, "should record two rolls")
		})
	})
}

func TestPlayer(t *testing.T) {
	t.Run("perfect_game", func(t *testing.T) {
		player := NewPlayer("max")
		// For frames 0-8, set as strikes.
		for i := 0; i < 9; i++ {
			nf, ok := player.Frames[i].(*NormalFrame)
			assert.True(t, ok, "frame %d should be a NormalFrame", i)
			err := nf.KnockPins(true, false)
			assert.NoError(t, err)
		}
		// Last frame as strike with two bonus rolls.
		lf, ok := player.Frames[9].(*LastFrame)
		assert.True(t, ok, "frame 9 should be a LastFrame")
		err := lf.KnockPins(true, false, 10, 10)
		assert.NoError(t, err)

		scores := player.GetScores()
		// In a perfect game, each frame should score 30.
		expected := []int{30, 30, 30, 30, 30, 30, 30, 30, 30, 30}
		assert.Equal(t, expected, scores, "perfect game should yield 300 total with each frame scoring 30")
	})

	t.Run("all_open_game", func(t *testing.T) {
		// setup
		player := NewPlayer("open")
		scores := [][]int{
			{3, 4}, // 7
			{2, 3}, // 5
			{4, 5}, // 9
			{3, 6}, // 9
			{2, 5}, // 7
			{4, 3}, // 7
			{1, 6}, // 7
			{2, 7}, // 9
			{3, 5}, // 8
		}
		for i, s := range scores {
			nf, ok := player.Frames[i].(*NormalFrame)
			require.True(t, ok, "frame %d should be a NormalFrame", i)
			err := nf.KnockPins(false, false, s[0], s[1])
			require.NoError(t, err)
		}
		// Last frame as open frame.
		lf, ok := player.Frames[9].(*LastFrame)
		require.True(t, ok, "frame 9 should be a LastFrame")
		err := lf.KnockPins(false, false, 3, 4) // 7
		require.NoError(t, err)

		// execute
		res := player.GetScores()

		// For open frames, score equals sum of frame's rolls.
		expected := []int{7, 5, 9, 9, 7, 7, 7, 9, 8, 7}
		assert.Equal(t, expected, res, "should return frame scores without bonuses")
	})

	t.Run("incomplete_game_with_strike_spare", func(t *testing.T) {
		player := NewPlayer("spare")
		player.Frames[0].(*NormalFrame).KnockPins(true, false)
		player.Frames[1].(*NormalFrame).KnockPins(false, true, 4)
		player.Frames[2].(*NormalFrame).KnockPins(false, false, 4, 4)
		player.Frames[3].(*NormalFrame).KnockPins(true, false)
		expected := []int{20, 14, 8, 10, 0, 0, 0, 0, 0, 0}
		assert.Equal(t, expected, player.GetScores())
	})
}
