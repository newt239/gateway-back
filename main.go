package main

import (
	"net/http"

	authRoute "github.com/newt239/gateway-back/routes/auth"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func Hello() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World")
	}
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	v1 := e.Group("/v1")
	v1.GET("/", Hello())

	auth := v1.Group("/auth")
	auth.POST("/login", authRoute.Login())
	auth.GET("/me", authRoute.Me())

	activity := v1.Group("/activity")
	activity.Use(middleware.JWT([]byte("secret")))

	e.Start(":3000")
}
