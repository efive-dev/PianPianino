package routes

import (
	"pianpianino/handlers"
	"pianpianino/helpers"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupRoutes(e *echo.Echo, auth *handlers.AuthHandler, task *handlers.TaskHandler) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "http://localhost:1323/"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Public routes using struct methods
	e.POST("/register", auth.Register)
	e.POST("/login", auth.Login)

	// Protected routes
	protected := e.Group("/api")
	jwtSecret := helpers.LoadConfig("JWT_SECRET")
	protected.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(jwtSecret),
		TokenLookup: "header:Authorization:Bearer ",
	}))

	protected.GET("/tasks", task.GetAllTasks)
	protected.POST("/tasks", task.InsertTask)
	protected.DELETE("/tasks/:id", task.DeleteTask)
	protected.PATCH("/tasks/:id/toggle", task.ToggleTaskCompleted)
}
