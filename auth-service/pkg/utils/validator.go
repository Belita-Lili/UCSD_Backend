package utils

import (
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var (
	validate   *validator.Validate
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

func init() {
	validate = validator.New()
	_ = validate.RegisterValidation("password", validatePassword)
}

// ValidateStruct valida una estructura usando tags
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// validatePassword valida que la contraseña cumpla con los requisitos
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	var (
		hasMinLen  = len(password) >= 8
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasMinLen && hasUpper && hasLower && (hasNumber || hasSpecial)
}

// ValidateEmail verifica si un email es válido
func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}
