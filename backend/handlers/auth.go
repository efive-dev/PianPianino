package handlers

import (
	"net/http"
	"pianpianino/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/bcrypt"
)

// Add a struct to bind the JSON request body
type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthHandler struct {
	DB        *bun.DB
	JWTSecret string
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req UserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request format"})
	}

	if req.Username == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Username and password are required"})
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error in hashing password"})
	}

	user := &models.User{
		Username: req.Username,
		Password: string(hashPassword),
	}

	_, err = h.DB.NewInsert().
		Model(user).
		Exec(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Error inserting user"})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "User registered successfully"})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req UserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request format"})
	}

	if req.Username == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Username and password are required"})
	}

	user := new(models.User)
	err := h.DB.NewSelect().
		Model(user).
		Where("username = ?", req.Username).
		Scan(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid credentials"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid credentials"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(2 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.JWTSecret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not create token"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Login successful",
		"token":   tokenString,
	})
}
