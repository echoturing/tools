package middlewaretools

import (
	"net/http"

	"github.com/echoturing/log"
	"github.com/labstack/echo"
)

type APICode int

const (
	APICodeOK  APICode = 0
	APICodeErr APICode = 500
)

type HandlerFunc func(c echo.Context) (data interface{}, err error)

func HandlerFuncWrapper(fn HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		type response struct {
			Code      int         `json:"code"`
			Message   string      `json:"msg"`
			RequestID string      `json:"request_id"`
			Data      interface{} `json:"data"`
		}
		resp := response{
			Code:    int(APICodeOK),
			Message: "",
		}
		ctx := c.Request().Context()
		requestIDUserID := log.FromContext(ctx)
		data, err := fn(c)
		responseContentType := c.Response().Header().Get(echo.HeaderContentType)
		if responseContentType == echo.MIMEOctetStream {
			// 如果是stream。那就直接返回nil了。。因为data肯定是nil
			return nil
		}
		if err != nil {
			resp.Code = int(APICodeErr)
			resp.Message = err.Error()
		}
		resp.Data = data
		resp.RequestID = requestIDUserID.RequestID
		return c.JSON(http.StatusOK, &resp)
	}
}
