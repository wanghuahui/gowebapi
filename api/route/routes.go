package route

import (
	"net/http"

	"github.com/labstack/echo"
)

// Load returns the routes and middleware
func Load(e *echo.Echo) {

	routes(e)
}

// *****************************************************************************
// Routes
// *****************************************************************************

func routes(e *echo.Echo) {
	// e.GET("/user", )
	e.GET("/hello", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\n")
	})
}
