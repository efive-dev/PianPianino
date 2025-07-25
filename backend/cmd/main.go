package main

import (
	"pianpianino/database"
	"pianpianino/models"
	"pianpianino/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	database.InitDB()
	models.Migrate()
	e := echo.New()
	routes.SetupRoutes(e)
	e.Logger.Fatal(e.Start(":1323"))
}
