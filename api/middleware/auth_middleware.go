package middleware

import (
	"KaungHtetHein116/IVY-backend/api/transport"
	"KaungHtetHein116/IVY-backend/pkg/constants"
	"net/http"
	"strings"

	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/clerk/clerk-sdk-go/v2/user"
	"github.com/labstack/echo/v4"
)

type paths struct {
	path   string
	method string
}

var openRoutes = []paths{
	{
		path:   "/",
		method: http.MethodGet,
	},
	{
		path:   "/api/v1/user/login",
		method: http.MethodPost,
	},
	{
		path:   "/api/v1/user/register",
		method: http.MethodPost,
	},

	{
		path:   "/api/v1/branch",
		method: http.MethodGet,
	},
	{
		path:   "/api/v1/branch/:id",
		method: http.MethodPut,
	},

	{
		path:   "/api/v1/category",
		method: http.MethodGet,
	},
	{
		path:   "/api/v1/category/:id",
		method: http.MethodGet,
	},

	{
		path:   "/api/v1/service",
		method: http.MethodGet,
	},
	{
		path:   "/api/v1/service/:id",
		method: http.MethodGet,
	},

	{
		path:   "/api/v1/booking",
		method: http.MethodGet,
	},
	{
		path:   "/api/v1/booking/:id",
		method: http.MethodGet,
	},
	{
		path:   "/api/v1/user/clerk-user-webhook",
		method: http.MethodPost,
	},
}

func RegisterAuthMiddleware(e *echo.Echo) {
	e.Use(JWTMiddleware())
}

func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if !shouldAuthorize(c) {
				return next(c)
			}

			var sessionToken string

			tokenCookie, err := c.Cookie("__session")

			if err == nil && tokenCookie.Value != "" {
				sessionToken = tokenCookie.Value
			} else if authHeader := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer "); authHeader != "" {
				sessionToken = authHeader
			}

			if sessionToken == "" {
				return transport.NewApiErrorResponse(c, http.StatusUnauthorized, constants.ErrUnauthorized, nil)
			}

			claims, err := jwt.Verify(c.Request().Context(), &jwt.VerifyParams{
				Token: sessionToken,
			})

			if err != nil {
				return transport.NewApiErrorResponse(c, http.StatusUnauthorized, constants.ErrInvalidToken, nil)
			}

			usr, err := user.Get(c.Request().Context(), claims.Subject)

			if err != nil {
				return transport.NewApiErrorResponse(c, http.StatusUnauthorized, constants.ErrInvalidToken, nil)
			}

			c.Set("user_id", usr.ID)

			return next(c)
		}
	}
}

func shouldAuthorize(c echo.Context) bool {
	requestPath := c.Path()
	requestMethod := c.Request().Method

	// Check if the current route is in excluded routes
	for _, route := range openRoutes {
		if route.path == requestPath && route.method == requestMethod {
			return false
		}
	}

	return true
}
