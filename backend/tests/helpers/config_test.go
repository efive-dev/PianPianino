package helpers_test

import (
	"os"
	"pianpianino/helpers"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigSuccess(t *testing.T) {
	envContent := "TEST_VAR=test_value\n"

	err := os.WriteFile(".env", []byte(envContent), 0644)
	assert.NoError(t, err)

	defer os.Remove(".env")

	result := helpers.LoadConfig("TEST_VAR")
	assert.Equal(t, "test_value", result)
}

func TestLoadConfigNonExistentVariable(t *testing.T) {
	envContent := "EXISTING_VAR=value\n"

	err := os.WriteFile(".env", []byte(envContent), 0644)
	assert.NoError(t, err)

	defer os.Remove(".env")

	result := helpers.LoadConfig("NON_EXISTENT")
	assert.Equal(t, "", result)
}

func TestLoadConfigMultipleVariables(t *testing.T) {
	envContent := `DB_HOST=localhost
		DB_PORT=5432
		JWT_SECRET=secret123
		`

	err := os.WriteFile(".env", []byte(envContent), 0644)
	assert.NoError(t, err)

	defer os.Remove(".env")

	assert.Equal(t, "localhost", helpers.LoadConfig("DB_HOST"))
	assert.Equal(t, "5432", helpers.LoadConfig("DB_PORT"))
	assert.Equal(t, "secret123", helpers.LoadConfig("JWT_SECRET"))
}

func TestLoadConfigEmptyValue(t *testing.T) {
	envContent := "EMPTY_VAR=\n"

	err := os.WriteFile(".env", []byte(envContent), 0644)
	assert.NoError(t, err)

	defer os.Remove(".env")

	result := helpers.LoadConfig("EMPTY_VAR")
	assert.Equal(t, "", result)
}
