package middleware

import (
	"KaungHtetHein116/IVY-backend/api/transport"
	"KaungHtetHein116/IVY-backend/pkg/constants"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type paths struct {
	path   string
	method string
}

var excludedRoutes = []paths{
	{
		path:   "/",
		method: http.MethodGet,
	},
	{
		path:   constants.USER_API_PREFIX + "/login",
		method: http.MethodPost,
	},
	{
		path:   constants.USER_API_PREFIX + "/register",
		method: http.MethodPost,
	},
}

func AuthMiddleware(e *echo.Echo) {
	e.Use(JWTMiddleware())
}

func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if !shouldAuthorize(c) {
				return next(c)
			}

			authHeader := c.Request().Header.Get("Authorization")

			if authHeader == "" {
				return transport.NewApiErrorResponse(c, http.StatusUnauthorized, constants.ErrUnauthorized, nil)
			}

			tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
			token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
				return []byte(os.Getenv("JWT_SECRET")), nil
			})

			if err != nil || !token.Valid {
				return transport.NewApiErrorResponse(c, http.StatusUnauthorized, constants.ErrInvalidToken, nil)
			}

			// Add user id to context
			claims := token.Claims.(jwt.MapClaims)
			c.Set("user_id", uint(claims["user_id"].(float64)))

			return next(c)
		}
	}
}

func shouldAuthorize(c echo.Context) bool {
	requestPath := c.Path()
	requestMethod := c.Request().Method

	// Check if the current route is in excluded routes
	for _, route := range excludedRoutes {
		if route.path == requestPath && route.method == requestMethod {
			return false
		}
	}
	return true
}
