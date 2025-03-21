package http_handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"bowling-score-tracker/configs"
	"bowling-score-tracker/core"
	"bowling-score-tracker/http_handlers/mocks"
)

func TestGameHttpHandler(t *testing.T) {
	t.Run("StartGame", func(t *testing.T) {
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
				mock.EXPECT().StartGame(gomock.Any(), gomock.Any()).Return(core.GameInfo{}, errors.New("abc"))

				// execute
				recorder := httptest.NewRecorder()
				req, _ := http.NewRequest(http.MethodPost, "/start", bytes.NewBuffer(body))
				r.ServeHTTP(recorder, req)

				// verify
				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			})
			t.Run("should_success_with_gameid_when_starting_game_successfully", func(t *testing.T) {
				// setup
				mock.EXPECT().StartGame(gomock.Any(), gomock.Any()).Return(core.GameInfo{Id: 4}, nil)

				// execute
				recorder := httptest.NewRecorder()
				req, _ := http.NewRequest(http.MethodPost, "/start", bytes.NewBuffer(body))
				r.ServeHTTP(recorder, req)

				// verify
				assert.Equal(t, http.StatusOK, recorder.Code)
				var response GameResponse
				require.Nil(t, json.Unmarshal(recorder.Body.Bytes(), &response))
				assert.Equal(t, int32(4), response.Id)
			})
		})
	})

	t.Run("GetGame", func(t *testing.T) {
		t.Run("should_return_bad_request_when_game_id_is_invalid", func(t *testing.T) {
			r := gin.Default()
			handler := NewGameHttpHandler(nil)

			r.GET("/:game_id/", handler.GetGame)

			req, _ := http.NewRequest(http.MethodGet, "/abc/", nil)
			recorder := httptest.NewRecorder()
			r.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusBadRequest, recorder.Code)
		})

		t.Run("should_return_error_when_manager_get_game_fails", func(t *testing.T) {
			r := gin.Default()
			mockCtrl := gomock.NewController(t)

			mockManager := mocks.NewMockGameManager(mockCtrl)
			handler := NewGameHttpHandler(mockManager)
			r.GET("/:game_id/", handler.GetGame)

			mockManager.EXPECT().GetGame(int32(456)).Return(core.GameInfo{}, errors.New("next frame error"))

			req, _ := http.NewRequest(http.MethodGet, "/456/", nil)
			recorder := httptest.NewRecorder()
			r.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusBadRequest, recorder.Code)
		})

		t.Run("should_return_current_frame_when_manager_next_frame_succeeds", func(t *testing.T) {
			r := gin.Default()
			mockCtrl := gomock.NewController(t)

			mockManager := mocks.NewMockGameManager(mockCtrl)
			handler := NewGameHttpHandler(mockManager)
			r.GET("/:game_id/", handler.GetGame)

			mockManager.EXPECT().GetGame(int32(789)).Return(core.GameInfo{CurrentFrame: 5}, nil)

			req, _ := http.NewRequest(http.MethodGet, "/789/", nil)
			recorder := httptest.NewRecorder()
			r.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response GameResponse
			require.Nil(t, json.Unmarshal(recorder.Body.Bytes(), &response))
			assert.Equal(t, 5, response.CurrentFrame)
		})
	})

	t.Run("SetFrameResult", func(t *testing.T) {
		t.Run("should_return_bad_request_when_input_data_cant_be_parsed", func(t *testing.T) {
			r := gin.Default()
			handler := NewGameHttpHandler(nil)

			r.POST("/:game_id/set_frame_result", handler.SetFrameResult)

			req, _ := http.NewRequest(http.MethodPost, "/123/set_frame_result", bytes.NewBuffer([]byte("abc")))
			recorder := httptest.NewRecorder()
			r.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusBadRequest, recorder.Code)
		})

		t.Run("should_return_bad_request_when_game_id_is_invalid", func(t *testing.T) {
			validReq := SetFrameResultRequest{
				PlayerIndex: 0,
				Pins:        []string{"X", "5"},
			}
			body, _ := json.Marshal(validReq)
			r := gin.Default()
			handler := NewGameHttpHandler(nil)

			r.POST("/:game_id/set_frame_result", handler.SetFrameResult)

			req, _ := http.NewRequest(http.MethodPost, "/abc/set_frame_result", bytes.NewBuffer(body))
			recorder := httptest.NewRecorder()
			r.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusBadRequest, recorder.Code)
		})

		t.Run("should_return_bad_request_when_parsing_pins_fails", func(t *testing.T) {
			r := gin.Default()
			handler := NewGameHttpHandler(nil)

			r.POST("/:game_id/set_frame_result", handler.SetFrameResult)

			invalidReq := SetFrameResultRequest{
				PlayerIndex: 0,
				Pins:        []string{"invalid"},
			}
			bodyInvalid, _ := json.Marshal(invalidReq)
			req, _ := http.NewRequest(http.MethodPost, "/123/set_frame_result", bytes.NewBuffer(bodyInvalid))
			recorder := httptest.NewRecorder()
			r.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusBadRequest, recorder.Code)
		})

		t.Run("when_input_is_valid", func(t *testing.T) {
			validReq := SetFrameResultRequest{
				PlayerIndex: 0,
				Pins:        []string{"X", "5"},
			}
			body, _ := json.Marshal(validReq)

			r := gin.Default()
			mockCtrl := gomock.NewController(t)

			mockManager := mocks.NewMockGameManager(mockCtrl)
			handler := NewGameHttpHandler(mockManager)
			r.POST("/:game_id/set_frame_result", handler.SetFrameResult)

			t.Run("should_call_manager_set_frame_result_with_correct_data", func(t *testing.T) {
				expectedPins := []interface{}{10, 5}

				mockManager.EXPECT().
					SetFrameResult(int32(123), validReq.PlayerIndex, expectedPins...).
					Return(core.GameInfo{}, nil)

				req, _ := http.NewRequest(http.MethodPost, "/123/set_frame_result", bytes.NewBuffer(body))
				recorder := httptest.NewRecorder()
				r.ServeHTTP(recorder, req)

				assert.Equal(t, http.StatusOK, recorder.Code)
			})

			t.Run("should_return_error_when_manager_set_frame_result_fails", func(t *testing.T) {
				mockManager.EXPECT().
					SetFrameResult(int32(123), validReq.PlayerIndex, gomock.Any()).
					Return(core.GameInfo{}, errors.New("set frame error"))

				req, _ := http.NewRequest(http.MethodPost, "/123/set_frame_result", bytes.NewBuffer(body))
				recorder := httptest.NewRecorder()
				r.ServeHTTP(recorder, req)

				assert.Equal(t, http.StatusBadRequest, recorder.Code)
			})
		})
	})

	t.Run("NextFrame", func(t *testing.T) {
		t.Run("should_return_bad_request_when_game_id_is_invalid", func(t *testing.T) {
			r := gin.Default()
			handler := NewGameHttpHandler(nil)

			r.POST("/:game_id/next_frame", handler.NextFrame)

			req, _ := http.NewRequest(http.MethodPost, "/abc/next_frame", nil)
			recorder := httptest.NewRecorder()
			r.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusBadRequest, recorder.Code)
		})

		t.Run("should_return_error_when_manager_next_frame_fails", func(t *testing.T) {
			r := gin.Default()
			mockCtrl := gomock.NewController(t)

			mockManager := mocks.NewMockGameManager(mockCtrl)
			handler := NewGameHttpHandler(mockManager)
			r.POST("/:game_id/next_frame", handler.NextFrame)

			mockManager.EXPECT().NextFrame(int32(456)).Return(core.GameInfo{}, errors.New("next frame error"))

			req, _ := http.NewRequest(http.MethodPost, "/456/next_frame", nil)
			recorder := httptest.NewRecorder()
			r.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusBadRequest, recorder.Code)
		})

		t.Run("should_return_current_frame_when_manager_next_frame_succeeds", func(t *testing.T) {
			r := gin.Default()
			mockCtrl := gomock.NewController(t)

			mockManager := mocks.NewMockGameManager(mockCtrl)
			handler := NewGameHttpHandler(mockManager)
			r.POST("/:game_id/next_frame", handler.NextFrame)

			mockManager.EXPECT().NextFrame(int32(789)).Return(core.GameInfo{CurrentFrame: 5}, nil)

			req, _ := http.NewRequest(http.MethodPost, "/789/next_frame", nil)
			recorder := httptest.NewRecorder()
			r.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusOK, recorder.Code)
			var response GameResponse
			require.Nil(t, json.Unmarshal(recorder.Body.Bytes(), &response))
			assert.Equal(t, 5, response.CurrentFrame)
		})
	})
}
