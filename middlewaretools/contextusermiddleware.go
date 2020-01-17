package middlewaretools

import (
	"encoding/json"

	"github.com/echoturing/log"
	"github.com/labstack/echo"
)

const HeaderUserID = "X-User-Id"

func WrapContextWithUser() func(echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			resp := c.Response()
			userID := req.Header.Get(HeaderUserID)
			requestID := resp.Header().Get(echo.HeaderXRequestID)
			ctx := c.Request().Context()
			requestIDUser := log.RequestIDWithUser{
				UserID:    json.Number(userID),
				RequestID: requestID,
			}
			ctx = log.NewContext(ctx, requestIDUser)
			c.SetRequest(req.WithContext(ctx))
			return next(c)
		}
	}
}
