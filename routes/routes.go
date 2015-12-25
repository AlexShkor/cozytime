package routes

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"

	"bitbucket.org/AlexShkor/cozytime/data"
	"bitbucket.org/AlexShkor/cozytime/models"
    "bitbucket.org/AlexShkor/cozytime/services"
    
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

func (r *Router) Authorize(c *echo.Context) error {
	decoder := json.NewDecoder(c.Request().Body)
	var model models.Authorization
	err := decoder.Decode(&model)
	if err != nil {
		return err
	}

	token, err := r.tokens.Authorize(model.PhoneNumber)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden)
	}

	return c.String(http.StatusOK, token)
}

func (r *Router) AdminIndex(c *echo.Context) error {
	html, err := getTemplate("views/index.html")
	if err != nil {
		return err
	}
	return c.HTML(http.StatusOK, html)
}

func (r *Router) SendSms(c *echo.Context) error {
	code, err := twilio.SendCode()
    if err != nil {
        c.String(http.StatusInternalServerError, "Sending faild")
    }
    phone := "+375259005003"
    setCodeToPhone(code, phone)
    return c.String(http.StatusOK, "Code sent!")
}

func setCodeToPhone(code string, phone string){
    return;
}
