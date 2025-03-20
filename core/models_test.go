package core

import "testing"
import "github.com/stretchr/testify/assert"

func TestTenPinGame(t *testing.T) {
	t.Run("when_start_game", func(t *testing.T) {
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
}
