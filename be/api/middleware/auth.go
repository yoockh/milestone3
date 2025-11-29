package middleware

import (
	"net/http"
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// for authentication middleware, such as JWT validation
func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		jwtSecretKey := os.Getenv("SECRET_KEY")
		jwtMiddleware := echojwt.WithConfig(echojwt.Config{
			SigningKey: []byte(jwtSecretKey),
			ErrorHandler: jwtErrorHandler,
		})

		return jwtMiddleware(next)(c)
	}
}

func jwtErrorHandler(c echo.Context, err error) error {
	return c.JSON(http.StatusUnauthorized, map[string]interface{}{
		"message": "you are unauthorized",
	})
}