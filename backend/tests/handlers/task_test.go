package handlers_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pianpianino/handlers"
	"pianpianino/models"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

const testJWTSecret = "test-secret-key"

func setUpTaskTestDB(t *testing.T) *bun.DB {
	// in memory db for testing
	ctx := context.Background()
	sqlDB, err := sql.Open(sqliteshim.ShimName, ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	DB := bun.NewDB(sqlDB, sqlitedialect.New())

	// migrate user table
	_, err = DB.NewCreateTable().
		Model((*models.User)(nil)).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// migrate task table
	_, err = DB.NewCreateTable().
		Model((*models.Task)(nil)).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		// Truncate tables after test runs
		_, err := DB.NewTruncateTable().
			Model((*models.Task)(nil)).
			Exec(ctx)
		if err != nil {
			t.Fatal(err)
		}
		_, err = DB.NewTruncateTable().
			Model((*models.User)(nil)).
			Exec(ctx)
		if err != nil {
			t.Fatal(err)
		}
		_ = sqlDB.Close()
	})

	return DB
}

func createTestJWTToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": float64(userID),
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString([]byte(testJWTSecret))
}

func createTestUser(t *testing.T, DB *bun.DB) int {
	user := &models.User{
		Username: "testuser",
		Password: "hashedpassword",
	}

	_, err := DB.NewInsert().
		Model(user).
		Exec(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	return int(user.ID)
}

func createTestTask(t *testing.T, DB *bun.DB, userID int, description string, priority models.Importance) *models.Task {
	task := &models.Task{
		UserID:      int64(userID),
		Description: description,
		Priority:    priority,
		Completed:   false,
	}

	_, err := DB.NewInsert().
		Model(task).
		Exec(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	return task
}

func TestGetAllTasksSuccess(t *testing.T) {
	DB := setUpTaskTestDB(t)
	handler := &handlers.TaskHandler{DB: DB}
	e := echo.New()

	userID := createTestUser(t, DB)

	createTestTask(t, DB, userID, "Task 1", models.Low)
	createTestTask(t, DB, userID, "Task 2", models.High)

	token, err := createTestJWTToken(userID)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	jwtToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(testJWTSecret), nil
	})
	ctx.Set("user", jwtToken)

	err = handler.GetAllTasks(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	tasks := response["tasks"].([]interface{})
	assert.Equal(t, 2, len(tasks))
	assert.Equal(t, float64(2), response["count"])
}

func TestGetAllTasksEmptyList(t *testing.T) {
	DB := setUpTaskTestDB(t)
	handler := &handlers.TaskHandler{DB: DB}
	e := echo.New()

	userID := createTestUser(t, DB)

	token, err := createTestJWTToken(userID)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	jwtToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(testJWTSecret), nil
	})
	ctx.Set("user", jwtToken)

	err = handler.GetAllTasks(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	tasks := response["tasks"].([]interface{})
	assert.Equal(t, 0, len(tasks))
	assert.Equal(t, float64(0), response["count"])
}

func TestInsertTaskSuccess(t *testing.T) {
	DB := setUpTaskTestDB(t)
	handler := &handlers.TaskHandler{DB: DB}
	e := echo.New()

	userID := createTestUser(t, DB)

	token, err := createTestJWTToken(userID)
	assert.NoError(t, err)

	reqBody := &handlers.TaskRequest{
		Description: "New test task",
		Priority:    models.Medium,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	jwtToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(testJWTSecret), nil
	})
	ctx.Set("user", jwtToken)

	err = handler.InsertTask(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Task created successfully", response["message"])
	assert.NotNil(t, response["task"])
}

func TestInsertTaskMissingDescription(t *testing.T) {
	DB := setUpTaskTestDB(t)
	handler := &handlers.TaskHandler{DB: DB}
	e := echo.New()

	userID := createTestUser(t, DB)

	token, err := createTestJWTToken(userID)
	assert.NoError(t, err)

	reqBody := &handlers.TaskRequest{
		Description: "",
		Priority:    models.Medium,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	jwtToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(testJWTSecret), nil
	})
	ctx.Set("user", jwtToken)

	err = handler.InsertTask(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestInsertTaskInvalidJSON(t *testing.T) {
	DB := setUpTaskTestDB(t)
	handler := &handlers.TaskHandler{DB: DB}
	e := echo.New()

	userID := createTestUser(t, DB)

	token, err := createTestJWTToken(userID)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader("{invalid json}"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	jwtToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(testJWTSecret), nil
	})
	ctx.Set("user", jwtToken)

	err = handler.InsertTask(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDeleteTaskSuccess(t *testing.T) {
	DB := setUpTaskTestDB(t)
	handler := &handlers.TaskHandler{DB: DB}
	e := echo.New()

	userID := createTestUser(t, DB)
	task := createTestTask(t, DB, userID, "Task to delete", models.Low)

	req := httptest.NewRequest(http.MethodDelete, "/tasks/"+strconv.Itoa(int(task.ID)), nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues(strconv.Itoa(int(task.ID)))

	err := handler.DeleteTask(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Task deleted successfully", response["message"])
}

func TestDeleteTaskInvalidID(t *testing.T) {
	DB := setUpTaskTestDB(t)
	handler := &handlers.TaskHandler{DB: DB}
	e := echo.New()

	req := httptest.NewRequest(http.MethodDelete, "/tasks/invalid", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues("invalid")

	err := handler.DeleteTask(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestToggleTaskCompletedSuccess(t *testing.T) {
	DB := setUpTaskTestDB(t)
	handler := &handlers.TaskHandler{DB: DB}
	e := echo.New()

	userID := createTestUser(t, DB)
	task := createTestTask(t, DB, userID, "Task to toggle", models.Medium)

	token, err := createTestJWTToken(userID)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPatch, "/tasks/"+strconv.Itoa(int(task.ID))+"/toggle", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues(strconv.Itoa(int(task.ID)))

	jwtToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(testJWTSecret), nil
	})
	ctx.Set("user", jwtToken)

	err = handler.ToggleTaskCompleted(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "Task completion toggled", response["message"])
}

func TestToggleTaskCompletedInvalidID(t *testing.T) {
	DB := setUpTaskTestDB(t)
	handler := &handlers.TaskHandler{DB: DB}
	e := echo.New()

	userID := createTestUser(t, DB)

	token, err := createTestJWTToken(userID)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPatch, "/tasks/invalid/toggle", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues("invalid")

	jwtToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(testJWTSecret), nil
	})
	ctx.Set("user", jwtToken)

	err = handler.ToggleTaskCompleted(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestToggleTaskCompletedTaskNotFound(t *testing.T) {
	DB := setUpTaskTestDB(t)
	handler := &handlers.TaskHandler{DB: DB}
	e := echo.New()

	userID := createTestUser(t, DB)

	token, err := createTestJWTToken(userID)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPatch, "/tasks/999/toggle", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues("999")

	jwtToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(testJWTSecret), nil
	})
	ctx.Set("user", jwtToken)

	err = handler.ToggleTaskCompleted(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestToggleTaskCompletedUserCannotAccessOtherUserTasks(t *testing.T) {
	DB := setUpTaskTestDB(t)
	handler := &handlers.TaskHandler{DB: DB}
	e := echo.New()

	userID1 := createTestUser(t, DB)

	user2 := &models.User{
		Username: "testuser2",
		Password: "hashedpassword2",
	}
	_, err := DB.NewInsert().
		Model(user2).
		Exec(context.Background())
	assert.NoError(t, err)
	userID2 := int(user2.ID)

	task := createTestTask(t, DB, userID1, "User1's task", models.High)

	token, err := createTestJWTToken(userID2)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPatch, "/tasks/"+strconv.Itoa(int(task.ID))+"/toggle", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("id")
	ctx.SetParamValues(strconv.Itoa(int(task.ID)))

	jwtToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(testJWTSecret), nil
	})
	ctx.Set("user", jwtToken)

	err = handler.ToggleTaskCompleted(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
