package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterBasicMiddlewares(e *echo.Echo) {
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "https://ivy-dashboard.vercel.app"},
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))
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
