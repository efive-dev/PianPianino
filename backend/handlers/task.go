package handlers

import (
	"net/http"
	"pianpianino/database"
	"pianpianino/models"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type TaskRequest struct {
	Description string            `json:"description" validate:"required"`
	Priority    models.Importance `json:"priority"`
}

// Helper function to get user ID from JWT token
func getUserIDFromToken(c echo.Context) (int, error) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(float64)
	return int(userID), nil
}

func GetAllTasks(c echo.Context) error {
	userID, err := getUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid token"})
	}

	DB := database.GetDB()
	tasks := make([]models.Task, 0)

	err = DB.NewSelect().
		Model(&tasks).
		Where("user_id = ?", userID).
		Scan(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch tasks"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"tasks": tasks,
		"count": len(tasks),
	})
}

func InsertTask(c echo.Context) error {
	userID, err := getUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid token"})
	}

	var req TaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}

	if req.Description == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Description is required"})
	}

	// Create new task
	task := models.Task{
		UserID:      int64(userID),
		Description: req.Description,
		Priority:    req.Priority,
		Completed:   false, // Default to false
	}

	DB := database.GetDB()

	_, err = DB.NewInsert().
		Model(&task).
		Exec(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create task"})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"message": "Task created successfully",
		"task":    task,
	})
}

func DeleteTask(c echo.Context) error {
	taskIDString := c.Param("id")
	taskID, err := strconv.Atoi(taskIDString)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "internal"})
	}

	task := new(models.Task)
	DB := database.GetDB()
	_, err = DB.NewDelete().
		Model(task).
		Where("id = ?", taskID).
		Exec(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "could not delete task"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "task deleted successfully"})
}

func ToggleTaskCompleted(c echo.Context) error {
	userID, err := getUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid token"})
	}

	taskIDStr := c.Param("id")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid task ID"})
	}

	DB := database.GetDB()
	task := new(models.Task)

	err = DB.NewSelect().
		Model(task).
		Where("id = ? AND user_id = ?", taskID, userID).
		Scan(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Task not found"})
	}

	task.Completed = !task.Completed

	_, err = DB.NewUpdate().
		Model(task).
		Column("completed").
		Where("id = ?", taskID).
		Exec(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to toggle task completion"})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Task completion toggled"})
}
