package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterBasicMiddlewares(e *echo.Echo) {
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `${time_rfc3339} ${status} ${method} ${uri} ${latency_human}` + "\n",
	}))
	e.Use(middleware.BodyLimit("1M"))

	e.GET("/", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, echo.Map{
			"message": "Welcome to IVY API",
		})
	})

}
