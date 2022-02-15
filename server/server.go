package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rbusquet/cosmic-go/entrypoints/allocate"
	"gorm.io/gorm"
)

func App(e *echo.Echo, db *gorm.DB) func() {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Use(middleware.Logger())
	e.Use(middleware.RemoveTrailingSlash())

	handler := allocate.Handler{DB: db}
	e.Use(handler.Transaction)
	e.POST("/allocate", handler.AllocateEndpoint)

	return func() {
		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}
}
