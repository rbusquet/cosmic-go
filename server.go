package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/rbusquet/cosmic-go/entrypoints/allocate"
	"github.com/rbusquet/cosmic-go/orm"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	db := orm.InitDB(&orm.Config{Debug: os.Getenv("DEBUG") == "1"})
	handler := allocate.Handler{DB: db}
	e.POST("/allocate", handler.AllocateEndpoint)
	e.Logger.Fatal(e.Start(":8080"))
}
