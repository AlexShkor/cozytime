package routes

import (
	"net/http"

	"bitbucket.org/AlexShkor/cozytime/data"

	"github.com/labstack/echo"
)

type GameRouter struct {
	games *data.Games
}

func NewGameRouter(games *data.Games) *GameRouter {
	return &GameRouter{games}
}

func (r *GameRouter) CreateGame(c *echo.Context) error {
	userId := getUserId(c)
	return c.JSON(http.StatusOK, userId)
}

func (r *GameRouter) JoinGame(c *echo.Context) error {
	userId := getUserId(c)

	return c.JSON(http.StatusOK, userId)
}

func (r *GameRouter) StartGame(c *echo.Context) error {
	userId := getUserId(c)

	return c.JSON(http.StatusOK, userId)
}

func (r *GameRouter) StopGame(c *echo.Context) error {
	userId := getUserId(c)

	return c.JSON(http.StatusOK, userId)
}

func getUserId(c *echo.Context) string {
	var data = c.Get("user")
	if userId, ok := data.(string); ok {
		return userId
	}
	return ""
}
