package main

import (
	"net/http"
	"pianpianino/database"
	"pianpianino/models"

	"github.com/labstack/echo/v4"
)

func main() {
	database.InitDB()
	models.Migrate()
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
