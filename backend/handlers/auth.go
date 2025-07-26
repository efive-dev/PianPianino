package handlers

import (
	"net/http"
	"pianpianino/database"
	"pianpianino/helpers"
	"pianpianino/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = helpers.LoadConfig("JWT_SECRET")

// Add a struct to bind the JSON request body
type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Register(c echo.Context) error {
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

	DB := database.GetDB()
	user := &models.User{
		Username: req.Username,
		Password: string(hashPassword),
	}

	_, err = DB.NewInsert().
		Model(user).
		Exec(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Error in inserting user in db"})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "User registered successfully"})
}

func Login(c echo.Context) error {
	var req UserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request format"})
	}

	if req.Username == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Username and password are required"})
	}

	DB := database.GetDB()
	user := new(models.User)
	err := DB.NewSelect().
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
		"iat":     time.Now().Unix(), // Issued at
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not create token"})
	}

	// Return both message and token
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Login successful",
		"token":   tokenString,
	})
}
