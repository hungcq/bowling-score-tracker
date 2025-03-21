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

				startGameRes, err := m.StartGame(configs.TenPin, []string{"hung"})

				assert.NoError(t, err)
				assert.GreaterOrEqual(t, startGameRes.Id, int32(1))
			})
		})
	})

	t.Run("GetGameInfo", func(t *testing.T) {
		t.Run("should_reject_invalid_game_id", func(t *testing.T) {
			m := NewGameManager()

			_, err := m.GetGame(1)

			assert.Error(t, err)
		})

		t.Run("should_return_game_info_when_game_id_is_valid", func(t *testing.T) {
			m := NewGameManager()
			startGameRes, err := m.StartGame(configs.TenPin, []string{"hung"})
			m.SetFrameResult(startGameRes.Id, 0, 10)

			res, err := m.GetGame(startGameRes.Id)

			assert.NoError(t, err)
			assert.Equal(t, GameInfo{
				Id:           startGameRes.Id,
				CurrentFrame: 0,
				Players: []PlayerScore{
					{
						Name:       "hung",
						Frames:     [][]int{{10}, nil, nil, nil, nil, nil, nil, nil, nil, nil},
						Scores:     []int{10, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						TotalScore: 10,
					},
				},
			}, res, "player score should reflect true score")
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
				startGameRes, _ := m.StartGame(configs.TenPin, []string{"hung"})

				_, err := m.SetFrameResult(startGameRes.Id, 0, 1)
				assert.Error(t, err)
			})
			t.Run("should_return_success_when_setting_frame_result_successfully", func(t *testing.T) {
				m := NewGameManager()
				startGameRes, err := m.StartGame(configs.TenPin, []string{"hung"})

				res, err := m.SetFrameResult(startGameRes.Id, 0, 10)

				assert.NoError(t, err)
				assert.Equal(t, GameInfo{
					Id:           startGameRes.Id,
					CurrentFrame: 0,
					Players: []PlayerScore{
						{
							Name:       "hung",
							Frames:     [][]int{{10}, nil, nil, nil, nil, nil, nil, nil, nil, nil},
							Scores:     []int{10, 0, 0, 0, 0, 0, 0, 0, 0, 0},
							TotalScore: 10,
						},
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
				startGameRes, err := m.StartGame(configs.TenPin, []string{"hung"})

				res, err := m.NextFrame(startGameRes.Id)

				assert.NoError(t, err)
				assert.Equal(t, 1, res.CurrentFrame)
			})
		})
	})
}
