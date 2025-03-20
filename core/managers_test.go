package core

import (
	"bowling-score-tracker/configs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGameManager(t *testing.T) {
	t.Run("when_start_game", func(t *testing.T) {
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
}
