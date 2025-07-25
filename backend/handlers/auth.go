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

func Register(c echo.Context) error {
	username := c.QueryParam("username")
	password := c.QueryParam("password")

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "error in hashing password"})
	}

	DB := database.GetDB()
	user := &models.User{
		Username: username,
		Password: string(hashPassword),
	}
	_, err = DB.NewInsert().
		Model(user).
		Exec(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "error in inserting user in db"})
	}

	return c.JSON(http.StatusCreated, echo.Map{"message": "user registered successfully"})
}

func Login(c echo.Context) error {
	username := c.QueryParam("username")
	password := c.QueryParam("password")

	DB := database.GetDB()
	user := new(models.User)
	err := DB.NewSelect().
		Model(user).
		Where("username = ?", username).
		Scan(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "error in retrieving user model"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid credentials"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"expiry":  time.Now().Add(2 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not create token"})
	}

	return c.JSON(http.StatusOK, echo.Map{"token": tokenString})
}
