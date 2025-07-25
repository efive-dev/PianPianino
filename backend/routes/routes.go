package routes

import (
	"pianpianino/handlers"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	// Public routes
	e.POST("/register", handlers.Register)
	e.POST("/login", handlers.Login)

}
