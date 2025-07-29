package main

import (
	"pianpianino/database"
	"pianpianino/handlers"
	"pianpianino/helpers"
	"pianpianino/models"
	"pianpianino/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	db := database.InitDB()
	models.Migrate()
	e := echo.New()

	authHandler := &handlers.AuthHandler{
		DB:        database.GetDB(),
		JWTSecret: helpers.LoadConfig("JWT_SECRET"),
	}

	taskHandler := &handlers.TaskHandler{DB: db}

	routes.SetupRoutes(e, authHandler, taskHandler)
	e.Logger.Fatal(e.Start(":1323"))
}
