package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"bitbucket.org/AlexShkor/cozytime/data"
	"bitbucket.org/AlexShkor/cozytime/models"
	"github.com/labstack/echo"
)

type GameRouter struct {
	games *data.Games
}

func NewGameRouter(games *data.Games) *GameRouter {
	return &GameRouter{games}
}

func (r *GameRouter) CreateGame(c *echo.Context) error {
	userID := getUserId(c)
	decoder := json.NewDecoder(c.Request().Body)
	var model models.CreateGame
	err := decoder.Decode(&model)
	fmt.Println("Model:")
	fmt.Println(model)
	if err != nil {
		return err
	}
	doc, err := r.games.Create(userID, model.Players, model.TargetTime)
	return c.JSON(http.StatusOK, models.GameResponse{doc.Id})
}

func (r *GameRouter) JoinGame(c *echo.Context) error {
	userID := getUserId(c)
	decoder := json.NewDecoder(c.Request().Body)
	var model models.JoinGame
	err := decoder.Decode(&model)
	fmt.Println("Model:")
	fmt.Println(model)
	if err != nil {
		return err
	}
	r.games.Join(model.GameId, userID)
	return c.JSON(http.StatusOK, userID)
}

func (r *GameRouter) StartGame(c *echo.Context) error {
	decoder := json.NewDecoder(c.Request().Body)
	var model models.JoinGame
	err := decoder.Decode(&model)
	fmt.Println("Model:")
	fmt.Println(model)
	if err != nil {
		return err
	}
	started, err := r.games.Start(model.GameId)
	return c.JSON(http.StatusOK, started)
}

func (r *GameRouter) StopGame(c *echo.Context) error {
	userID := getUserId(c)
	decoder := json.NewDecoder(c.Request().Body)
	var model models.JoinGame
	err := decoder.Decode(&model)
	fmt.Println("Model:")
	fmt.Println(model)
	if err != nil {
		return err
	}
	r.games.Stop(model.GameId, userID)
	return c.JSON(http.StatusOK, userID)
}

func getUserId(c *echo.Context) string {
	var data = c.Get("user")
	if userId, ok := data.(string); ok {
		return userId
	}
	return ""
}
