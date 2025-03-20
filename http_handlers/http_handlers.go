package http_handlers

import (
	"bowling-score-tracker/configs"
	"bowling-score-tracker/core"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterEndpoints(r *gin.Engine) {
	gameHandler := NewGameHttpHandler(core.NewGameManager())
	r.POST("/start_game", gameHandler.StartGame)
}

type GameHttpHandler struct {
	manager GameManager
}

func NewGameHttpHandler(manager GameManager) *GameHttpHandler {
	return &GameHttpHandler{
		manager: manager,
	}
}

//go:generate mockgen -source=http_handlers.go -destination=mocks/http_handlers.go -package=mocks
type GameManager interface {
	StartGame(t configs.GameType, playerNames []string) (gameId int32, err error)
}

type StartGameRequest struct {
	GameType    configs.GameType `json:"game_type"`
	PlayerNames []string         `json:"player_names" binding:"required,dive,max=5"`
}

type StartGameResponse struct {
	GameId int32 `json:"game_id"`
	Response
}

type Response struct {
	Error string `json:"error"`
}

func (h *GameHttpHandler) StartGame(c *gin.Context) {
	var req StartGameRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, StartGameResponse{
			Response: Response{
				Error: err.Error(),
			},
		})
		return
	}

	gameId, err := h.manager.StartGame(req.GameType, req.PlayerNames)
	if err != nil {
		c.JSON(http.StatusBadRequest, StartGameResponse{
			Response: Response{
				Error: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, &StartGameResponse{GameId: gameId})
}
