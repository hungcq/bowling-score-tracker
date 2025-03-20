package http_handlers

import (
	"bowling-score-tracker/configs"
	"bowling-score-tracker/http_handlers/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGameHttpHandler(t *testing.T) {
	t.Run("start_game", func(t *testing.T) {
		t.Run("should_return_bad_request_when_input_data_cant_be_parsed", func(t *testing.T) {
			r := gin.Default()
			handler := NewGameHttpHandler(nil)
			r.POST("/start", handler.StartGame)

			req, _ := http.NewRequest(http.MethodPost, "/start", bytes.NewBuffer([]byte("abc")))
			recorder := httptest.NewRecorder()
			r.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusBadRequest, recorder.Code)
		})
		t.Run("should_return_bad_request_when_request_is_invalid", func(t *testing.T) {
			r := gin.Default()
			handler := NewGameHttpHandler(nil)
			r.POST("/start", handler.StartGame)

			data := StartGameRequest{
				GameType:    configs.TenPin,
				PlayerNames: nil,
			}
			body, _ := json.Marshal(data)
			req, _ := http.NewRequest(http.MethodPost, "/start", bytes.NewBuffer(body))
			recorder := httptest.NewRecorder()
			r.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusBadRequest, recorder.Code)
		})
		t.Run("when_input_is_valid", func(t *testing.T) {
			// setup
			r := gin.Default()
			data := StartGameRequest{
				GameType:    configs.TenPin,
				PlayerNames: []string{"hung"},
			}
			body, _ := json.Marshal(data)
			mock := mocks.NewMockGameManager(gomock.NewController(t))
			handler := NewGameHttpHandler(mock)
			r.POST("/start", handler.StartGame)

			t.Run("should_start_game_with_correct_data", func(t *testing.T) {
				// verify
				mock.EXPECT().StartGame(data.GameType, data.PlayerNames).Times(1)

				// execute
				recorder := httptest.NewRecorder()
				req, _ := http.NewRequest(http.MethodPost, "/start", bytes.NewBuffer(body))
				r.ServeHTTP(recorder, req)
			})
			t.Run("should_return_error_when_failing_to_start_game", func(t *testing.T) {
				// setup
				mock.EXPECT().StartGame(gomock.Any(), gomock.Any()).Return(int32(0), errors.New("abc"))

				// execute
				recorder := httptest.NewRecorder()
				req, _ := http.NewRequest(http.MethodPost, "/start", bytes.NewBuffer(body))
				r.ServeHTTP(recorder, req)

				// verify
				assert.Equal(t, http.StatusInternalServerError, recorder.Code)
			})
			t.Run("should_success_with_gameid_when_starting_game_successfully", func(t *testing.T) {
				// setup
				mock.EXPECT().StartGame(gomock.Any(), gomock.Any()).Return(int32(4), nil)

				// execute
				recorder := httptest.NewRecorder()
				req, _ := http.NewRequest(http.MethodPost, "/start", bytes.NewBuffer(body))
				r.ServeHTTP(recorder, req)

				// verify
				assert.Equal(t, http.StatusOK, recorder.Code)
				var response StartGameResponse
				require.Nil(t, json.Unmarshal(recorder.Body.Bytes(), &response))
				assert.Equal(t, int32(4), response.GameId)
			})
		})
	})
}
