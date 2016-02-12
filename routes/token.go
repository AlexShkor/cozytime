package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"bitbucket.org/AlexShkor/cozytime/data"
	"bitbucket.org/AlexShkor/cozytime/models"
	"bitbucket.org/AlexShkor/cozytime/services"

	"github.com/astaxie/beego/validation"
	"github.com/labstack/echo"
)

func getTemplate(path string) (string, error) {
	var htmlBuffer bytes.Buffer
	tpl, err := template.ParseFiles(path)
	if err != nil {
		return "", err
	}
	tpl.Execute(&htmlBuffer, nil)
	return htmlBuffer.String(), nil
}

type Router struct {
	tokens *data.Tokens
}

func NewRouter(tokens *data.Tokens) *Router {
	return &Router{tokens}
}

func (r *Router) HelloWorld(c *echo.Context) error {
	docs, _ := r.tokens.GetAll()
	return c.JSON(http.StatusOK, docs)
}

func (r *Router) SubmitCode(c *echo.Context) error {
	fmt.Println("Authorize action!")
	decoder := json.NewDecoder(c.Request().Body)
	var model models.Authorization
	err := decoder.Decode(&model)
	fmt.Println("Model:")
	fmt.Println(model)
	if err != nil {
		return err
	}

	user, err := r.tokens.Authorize(model.PhoneNumber, model.Code)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Can't authorize")
	}
	return c.JSON(http.StatusOK, user)
}

func (r *Router) AdminIndex(c *echo.Context) error {
	html, err := getTemplate("views/index.html")
	if err != nil {
		return err
	}
	return c.HTML(http.StatusOK, html)
}

func (r *Router) SetName(c *echo.Context) error {
	decoder := json.NewDecoder(c.Request().Body)
	var model models.SetName
	err := decoder.Decode(&model)
	if err != nil {
		return err
	}
	var data = c.Get("user")
	if userId, ok := data.(string); ok {
		fmt.Println("User:")
		fmt.Println(userId)
		err := r.tokens.SetName(userId, model.Name)
		if err != nil {
			fmt.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Can't set name")
		}
	}
	return c.String(http.StatusOK, "Code sent!")
}

func (r *Router) IsAuthorized(c *echo.Context) error {
	return c.String(http.StatusOK, "Authorized!")
}

func (r *Router) Register(c *echo.Context) error {
	decoder := json.NewDecoder(c.Request().Body)
	var model models.Authorization
	err := decoder.Decode(&model)
	fmt.Println(model.PhoneNumber)
	valid := validation.Validation{}
	valid.Numeric(model.PhoneNumber, "phone")
	if err != nil || valid.HasErrors() {
		return echo.NewHTTPError(http.StatusInternalServerError, "Invalid phone.")
	}

	user, err := r.tokens.FindByPhone(model.PhoneNumber)
	fmt.Println("User:")
	fmt.Println(user)
	code := ""
	if user != nil {
		fmt.Println("User updating")
		code, err = r.tokens.UpdateCode(user.Id, model.PhoneNumber)
		if err != nil {
			fmt.Println(err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Can't update user info with new code.")
		}
	} else {
		fmt.Println("User creating")
		code, err = r.tokens.Create(model.PhoneNumber)
		if err != nil {
			fmt.Println(err)
			return echo.NewHTTPError(http.StatusForbidden)
		}
	}
	if code == "" {
		return echo.NewHTTPError(http.StatusInternalServerError, "Can't update user info with new code.")
	}
	err = twilio.SendCode(code, model.PhoneNumber)
	return c.String(http.StatusOK, "Code sent.")
}

func (r *Router) GetFriendsList(c *echo.Context) error {
	decoder := json.NewDecoder(c.Request().Body)
	var model models.Friends
	err := decoder.Decode(&model)
	if err != nil {
		return err
	}

	friends, err := r.tokens.FindFriends(model.PhoneNumbers)
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusForbidden)
	}

	return c.JSON(http.StatusOK, friends)
}
