package http_handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"bowling-score-tracker/configs"
	"bowling-score-tracker/core"
)

func RegisterEndpoints(r *gin.Engine) {
	gameHandler := NewGameHttpHandler(core.NewGameManager())
	r.POST("/start_game", gameHandler.StartGame)
	r.GET("/:game_id", gameHandler.GetGame)
	r.POST("/:game_id/set_frame_result", gameHandler.SetFrameResult)
	r.POST("/:game_id/next_frame", gameHandler.NextFrame)
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
	StartGame(t configs.GameType, playerNames []string) (core.GameInfo, error)
	GetGame(gameId int32) (core.GameInfo, error)
	SetFrameResult(gameId int32, playerIndex int, pins ...int) (core.GameInfo, error)
	NextFrame(gameId int32) (core.GameInfo, error)
}

type StartGameRequest struct {
	GameType    configs.GameType `json:"game_type"`
	PlayerNames []string         `json:"player_names" binding:"required,dive,max=5"`
}

type Response struct {
	Error string `json:"error,omitempty"`
}

func (h *GameHttpHandler) StartGame(c *gin.Context) {
	var req StartGameRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, GameResponse{
			Response: Response{
				Error: err.Error(),
			},
		})
		return
	}

	res, err := h.manager.StartGame(req.GameType, req.PlayerNames)
	if err != nil {
		c.JSON(http.StatusBadRequest, GameResponse{
			Response: Response{
				Error: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, GameResponse{GameInfo: &res})
}

type GameResponse struct {
	*core.GameInfo `json:"game,omitempty"`
	Response
}

func (h *GameHttpHandler) GetGame(c *gin.Context) {
	gameId, err := parseGameId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, GameResponse{
			Response: Response{
				Error: err.Error(),
			},
		})
		return
	}

	res, err := h.manager.GetGame(gameId)
	if err != nil {
		c.JSON(http.StatusBadRequest, GameResponse{
			Response: Response{
				Error: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, GameResponse{
		GameInfo: &res,
	})
}

type SetFrameResultRequest struct {
	PlayerIndex int      `json:"player_index" binding:"min=0"`
	Pins        []string `json:"pins" binding:"required,dive,min=1"`
}

func (h *GameHttpHandler) SetFrameResult(c *gin.Context) {
	var req SetFrameResultRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, GameResponse{
			Response: Response{
				Error: err.Error(),
			},
		})
		return
	}

	gameId, err := parseGameId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, GameResponse{
			Response: Response{
				Error: err.Error(),
			},
		})
		return
	}

	pins, err := parsePins(req.Pins)
	if err != nil {
		c.JSON(http.StatusBadRequest, GameResponse{
			Response: Response{
				Error: err.Error(),
			},
		})
	}

	res, err := h.manager.SetFrameResult(gameId, req.PlayerIndex, pins...)
	if err != nil {
		c.JSON(http.StatusBadRequest, GameResponse{
			Response: Response{
				Error: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, GameResponse{
		GameInfo: &res,
	})
}

func parsePins(pins []string) ([]int, error) {
	var res []int
	for i, str := range pins {
		pin, err := parsePin(str)
		if err != nil {
			return nil, err
		}
		if pin == spare {
			if i != 1 {
				return nil, errors.New("/ must be at index 1")
			}
			res = append(res, 10-res[i-1])
		} else {
			res = append(res, pin)
		}
	}
	return res, nil
}

const spare = -1

func parsePin(pin string) (int, error) {
	switch pin {
	case "X":
		return 10, nil
	case "/":
		return spare, nil
	case "-":
		return 0, nil
	default:
		i, err := strconv.Atoi(pin)
		if err != nil {
			return 0, err
		}
		if i < 0 || i >= 10 {
			return 0, errors.New("pin must be X, /, or between 0 and 10")
		}
		return i, nil
	}
}

func (h *GameHttpHandler) NextFrame(c *gin.Context) {
	gameId, err := parseGameId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, GameResponse{
			Response: Response{
				Error: err.Error(),
			},
		})
		return
	}

	res, err := h.manager.NextFrame(gameId)
	if err != nil {
		c.JSON(http.StatusBadRequest, GameResponse{
			Response: Response{
				Error: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, GameResponse{
		GameInfo: &res,
	})
}

func parseGameId(c *gin.Context) (int32, error) {
	// Get the "id" parameter from the path.
	idParam := c.Param("game_id")

	// Parse the parameter as int32.
	id64, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		return 0, errors.New("invalid id parameter")
	}
	id := int32(id64)
	return id, nil
}
