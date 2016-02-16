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
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Can't create game")
	}
	return c.JSON(http.StatusOK, models.GameResponse{doc.Id})
}

func (r *GameRouter) JoinGame(c *echo.Context) error {
	fmt.Println("JOIN GAME:")
	userID := getUserId(c)
	decoder := json.NewDecoder(c.Request().Body)
	var model models.JoinGame
	err := decoder.Decode(&model)
	err = r.games.Join(model.GameId, userID)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Can't join game")
	}
	return c.JSON(http.StatusOK, userID)
}

func (r *GameRouter) LeaveGame(c *echo.Context) error {
	fmt.Println("LEAVE GAME:")
	userID := getUserId(c)
	decoder := json.NewDecoder(c.Request().Body)
	var model models.JoinGame
	err := decoder.Decode(&model)
	err = r.games.Leave(model.GameId, userID)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Can't leave game")
	}
	return c.JSON(http.StatusOK, userID)
}

func (r *GameRouter) DeleteGame(c *echo.Context) error {
	fmt.Println("DELETE GAME:")
	userID := getUserId(c)
	decoder := json.NewDecoder(c.Request().Body)
	var model models.JoinGame
	err := decoder.Decode(&model)
	err = r.games.Delete(model.GameId, userID)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Can't delete game")
	}
	return c.JSON(http.StatusOK, userID)
}

func (r *GameRouter) StartGame(c *echo.Context) error {
	fmt.Println("START GAME:")
	decoder := json.NewDecoder(c.Request().Body)
	var model models.JoinGame
	err := decoder.Decode(&model)
	if err != nil {
		return err
	}
	started, err := r.games.Start(model.GameId)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Can't start game")
	}
	fmt.Println("START STARTED: ")
	fmt.Print(started)
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
	err = r.games.Stop(model.GameId, userID)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Can't stop game")
	}
	return c.JSON(http.StatusOK, userID)
}

func (r *GameRouter) GetGame(c *echo.Context) error {
	decoder := json.NewDecoder(c.Request().Body)
	var model models.JoinGame
	decoder.Decode(&model)
	doc, err := r.games.Get(model.GameId)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Can't get game")
	}
	return c.JSON(http.StatusOK, doc)
}

func (r *GameRouter) GetMyGames(c *echo.Context) error {
	userID := getUserId(c)
	docs, err := r.games.GetAllForUser(userID)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Can't get game")
	}
	return c.JSON(http.StatusOK, docs)
}

func getUserId(c *echo.Context) string {
	var data = c.Get("user")
	if userId, ok := data.(string); ok {
		return userId
	}
	return ""
}
