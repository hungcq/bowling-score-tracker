package core

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"bowling-score-tracker/configs"
)

func TestGameManager(t *testing.T) {
	t.Run("StartGame", func(t *testing.T) {
		t.Run("should_reject_invalid_game_type", func(t *testing.T) {
			m := NewGameManager()

			_, err := m.StartGame("abc", []string{"hung"})

			assert.Error(t, err)
		})

		t.Run("when_game_type_is_valid", func(t *testing.T) {
			t.Run("should_return_error_when_failing_to_start_game", func(t *testing.T) {
				m := NewGameManager()

				_, err := m.StartGame(configs.TenPin, []string{""})

				assert.Error(t, err)
			})
			t.Run("should_return_success_when_starting_game_successfully", func(t *testing.T) {
				m := NewGameManager()

				gameId, err := m.StartGame(configs.TenPin, []string{"hung"})

				assert.NoError(t, err)
				assert.GreaterOrEqual(t, gameId, int32(1))
			})
		})
	})

	t.Run("SetFrameResult", func(t *testing.T) {
		t.Run("should_reject_invalid_game_id", func(t *testing.T) {
			m := NewGameManager()

			_, err := m.SetFrameResult(1, 0, 1)

			assert.Error(t, err)
		})

		t.Run("when_game_id_is_valid", func(t *testing.T) {
			t.Run("should_return_error_when_failing_to_set_frame_result", func(t *testing.T) {
				m := NewGameManager()
				gameId, _ := m.StartGame(configs.TenPin, []string{"hung"})

				_, err := m.SetFrameResult(gameId, 0, 1)
				assert.Error(t, err)
			})
			t.Run("should_return_success_when_setting_frame_result_successfully", func(t *testing.T) {
				m := NewGameManager()
				gameId, err := m.StartGame(configs.TenPin, []string{"hung"})

				res, err := m.SetFrameResult(gameId, 0, 10)

				assert.NoError(t, err)
				assert.Equal(t, []PlayerScore{
					{
						Scores:     []int{10, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						TotalScore: 10,
					},
				}, res, "player score should reflect true score")
			})
		})
	})

	t.Run("NextFrame", func(t *testing.T) {
		t.Run("should_reject_invalid_game_id", func(t *testing.T) {
			m := NewGameManager()

			_, err := m.NextFrame(1)

			assert.Error(t, err)
		})

		t.Run("when_game_id_is_valid", func(t *testing.T) {
			t.Run("should_increment_frame", func(t *testing.T) {
				m := NewGameManager()
				gameId, err := m.StartGame(configs.TenPin, []string{"hung"})

				res, err := m.NextFrame(gameId)

				assert.NoError(t, err)
				assert.Equal(t, 1, res)
			})
		})
	})
}
