package main

import (
	"net/http"
	"strings"

	"bitbucket.org/AlexShkor/cozytime/data"

	"github.com/labstack/echo"
)

const (
	Bearer = "Bearer"
)

func BearerAuth(tokens *data.Tokens) echo.HandlerFunc {
	return func(c *echo.Context) error {
		unauthorizedError := echo.NewHTTPError(http.StatusUnauthorized)

		authHeader := c.Request().Header.Get("Authorization")

		if !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			return unauthorizedError
		}

		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 {
			return unauthorizedError
		}

		bearerToken := authParts[1]

		var userId string
		if userId, err := tokens.IsAuthorized(bearerToken); err != nil || userId == "" {
			return unauthorizedError
		}

		c.Set("user", userId)

		return nil
	}
}

func stripTrailingSlash(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	l := len(path) - 1
	if path != "/" && path[l] == '/' {
		r.URL.Path = path[:l]
	}
}
