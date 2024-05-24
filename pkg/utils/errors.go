package utils

import (
	"errors"
	"strings"

	"github.com/labstack/echo/v4"
)

// ErrInvalidURL is returned when URL is invalid.
var ErrInvalidURL = errors.New(
	"invalid URL. It won't be parsed. Check that your url contains scheme",
)

func Auth(ctx echo.Context) bool {
	const authToken = "dXNlcm5hbWU6cGFzc3dvcmQ=" // username:password
	var authorization = strings.Split(ctx.Request().Header.Get("Authorization"), "Basic")
	if (len(authorization) <= 1) {
		return false
	}
	return strings.TrimSpace(authorization[1]) == authToken;

}