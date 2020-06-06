package middlewaretools

import (
	"regexp"
	"sync"

	"github.com/echoturing/log"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	accessLogFormat = `{"level":"info","msg":"access_log","time":"${time_rfc3339_nano}","x-request-id":"${id}","remote_ip":"${remote_ip}",` +
		`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
		`"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}"` +
		`,"bytes_in":${bytes_in},"bytes_out":${bytes_out},"protocol":"${protocol}"` +
		`,"referer":"${referer}"}` + "\n"
)

var (
	passwordReg = regexp.MustCompile(`("password"\s*:\s*)"([a-zA-Z0-9_]+)"`)

	FilterPasswordFunc = func(s []byte) []byte {
		return passwordReg.ReplaceAll(s, []byte(`$1"*******"`))
	}
	BodyDumpMiddleware = func(filterRequestFuncs ...func([]byte) []byte) echo.MiddlewareFunc {
		return middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{
			Handler: func(c echo.Context, reqData []byte, respData []byte) {
				for _, f := range filterRequestFuncs {
					reqData = f(reqData)
				}
				log.InfoWithContext(c.Request().Context(), "body_access_log", "path", c.Request().URL, "req", string(reqData), "resp", string(respData))
			},
			Skipper: middleware.DefaultSkipper})
	}
	once   sync.Once
	config = middleware.DefaultLoggerConfig
)

func AccessLog(skipper middleware.Skipper) echo.MiddlewareFunc {
	once.Do(func() {
		config.Format = accessLogFormat
		if skipper != nil {
			config.Skipper = skipper
		}
	})
	return middleware.LoggerWithConfig(config)
}
