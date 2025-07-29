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
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

func setUpTestDB(t *testing.T) *bun.DB {
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

	t.Cleanup(func() {
		// Truncate the table after test runs
		_, err := DB.NewTruncateTable().
			Model((*models.User)(nil)).
			Exec(ctx)
		if err != nil {
			t.Fatal(err)
		}
		_ = sqlDB.Close()
	})

	return DB
}

func TestRegisterSuccess(t *testing.T) {
	DB := setUpTestDB(t)
	handler := &handlers.AuthHandler{DB: DB}
	e := echo.New()

	reqBody := &handlers.UserRequest{
		Username: "correctUser",
		Password: "correctPassword",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)

	err := handler.Register(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

func TestRegisterInvalidJSON(t *testing.T) {
	DB := setUpTestDB(t)
	handler := &handlers.AuthHandler{DB: DB}
	e := echo.New()

	req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader("{invalid json}"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)

	err := handler.Register(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestRegisterMissingFieldsUsername(t *testing.T) {
	DB := setUpTestDB(t)
	handler := &handlers.AuthHandler{DB: DB}
	e := echo.New()

	reqBody := &handlers.UserRequest{
		Username: "",
		Password: "testPassword",
	}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)

	err := handler.Register(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestRegisterMissingFieldsPassword(t *testing.T) {
	DB := setUpTestDB(t)
	handler := &handlers.AuthHandler{DB: DB}
	e := echo.New()

	reqBody := &handlers.UserRequest{
		Username: "testUsername",
		Password: "",
	}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)

	err := handler.Register(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestLoginSuccess(t *testing.T) {
	DB := setUpTestDB(t)
	handler := &handlers.AuthHandler{
		DB:        DB,
		JWTSecret: "supersecret",
	}
	e := echo.New()

	registerBody := &handlers.UserRequest{
		Username: "integrationUser",
		Password: "integrationPass",
	}
	jsonReg, _ := json.Marshal(registerBody)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonReg))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := handler.Register(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	loginBody := &handlers.UserRequest{
		Username: "integrationUser",
		Password: "integrationPass",
	}
	jsonLogin, _ := json.Marshal(loginBody)

	loginReq := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonLogin))
	loginReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	loginRec := httptest.NewRecorder()
	loginCtx := e.NewContext(loginReq, loginRec)

	err = handler.Login(loginCtx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, loginRec.Code)

	var loginResp map[string]interface{}
	err = json.Unmarshal(loginRec.Body.Bytes(), &loginResp)
	assert.NoError(t, err)

	assert.Equal(t, "Login successful", loginResp["message"])
	assert.NotEmpty(t, loginResp["token"])
}

func TestLoginWrongPassword(t *testing.T) {
	DB := setUpTestDB(t)
	handler := &handlers.AuthHandler{DB: DB, JWTSecret: "test"}
	e := echo.New()

	handler.Register(e.NewContext(
		httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(`{"username":"foo","password":"bar"}`)),
		httptest.NewRecorder(),
	))

	loginBody := &handlers.UserRequest{
		Username: "foo",
		Password: "wrongpass",
	}
	jsonBody, _ := json.Marshal(loginBody)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := handler.Login(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestLoginNonExistentUser(t *testing.T) {
	DB := setUpTestDB(t)
	handler := &handlers.AuthHandler{DB: DB, JWTSecret: "test"}
	e := echo.New()

	loginBody := &handlers.UserRequest{
		Username: "ghost",
		Password: "password",
	}
	jsonBody, _ := json.Marshal(loginBody)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := handler.Login(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestLoginInvalidJSON(t *testing.T) {
	DB := setUpTestDB(t)
	handler := &handlers.AuthHandler{DB: DB, JWTSecret: "test"}
	e := echo.New()

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("{bad json"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	err := handler.Login(ctx)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}
