package route

import (
	"gowebapi/api/controller"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Load returns the routes and middleware
func Load(e *echo.Echo) {
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	routes(e)
}

// *****************************************************************************
// Routes
// *****************************************************************************

func routes(e *echo.Echo) {
	// 接口版本信息
	e.GET("/version", func(c echo.Context) error {
		return c.String(http.StatusOK, "The Version is v1.0.0\n")
	})

	// 游客可以访问的接口
	// 用户注册
	e.POST("/users", controller.UserCreate)

	// 需要 token 验证的接口
	// Restricted group
	r := e.Group("")
	r.Use(middleware.JWT([]byte("secret")))
	// 当前登录用户信息
	r.GET("/user", restricted)
}
