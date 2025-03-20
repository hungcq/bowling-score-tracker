package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGameManager(t *testing.T) {
	t.Run("when_start_game", func(t *testing.T) {
		t.Run("should_reject_invalid_game_type", func(t *testing.T) {
			m := NewGameManager()
			assert.Error(t, m.StartGame("abc", []string{"hung"}))
		})

		t.Run("when_game_type_is_valid", func(t *testing.T) {
			t.Run("should_return_error_when_failing_to_start_game", func(t *testing.T) {
				m := NewGameManager()
				assert.Error(t, m.StartGame(TenPin, []string{""}))
			})
			t.Run("should_return_success_when_starting_game_successfully", func(t *testing.T) {
				m := NewGameManager()
				assert.NoError(t, m.StartGame(TenPin, []string{"hung"}))
			})
		})
	})
}
