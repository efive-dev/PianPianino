package routes

import (
	"pianpianino/handlers"
	"pianpianino/helpers"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupRoutes(e *echo.Echo) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "http://localhost:1323/"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))
	// Public routes
	e.POST("/register", handlers.Register)
	e.POST("/login", handlers.Login)

	protected := e.Group("/api")
	jwtSecret := helpers.LoadConfig("JWT_SECRET")
	protected.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(jwtSecret),
		TokenLookup: "header:Authorization:Bearer ",
	}))

	protected.GET("/tasks", handlers.GetAllTasks)
	protected.POST("/tasks", handlers.InsertTask)
	protected.DELETE("/tasks/:id", handlers.DeleteTask)
	protected.PATCH("/tasks/:id/toggle", handlers.ToggleTaskCompleted)

}
