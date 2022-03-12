package helpers

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ValidateUser(c *gin.Context, dataSet interface{}) (bool, interface{}) {
	if err := c.ShouldBindJSON(dataSet); err != nil {
		if _, ok := err.(validator.ValidationErrors); !ok {
			return false, err.Error()
		}

		errors := make(gin.H)

		for _, err := range err.(validator.ValidationErrors) {
			name := strings.ToLower(err.StructField())
			switch err.Tag() {
			case "required":
				errors[name] = name + " is required"
			case "email":
				errors[name] = name + " should be a valid email"
			case "min":
				errors[name] = name + " minimum " + err.Param() + " characters"
			case "max":
				errors[name] = name + " maximum " + err.Param() + " characters"
			default:
				errors[name] = name + " is invalid"
			}
		}
		return false, errors
	}
	return true, nil

}

func ValidateServe(c *gin.Context, dataSet interface{}) (bool, interface{}) {
	if err := c.ShouldBindJSON(dataSet); err != nil {
		if _, ok := err.(validator.ValidationErrors); !ok {
			return false, err.Error()
		}

		errors := make(gin.H)

		for _, err := range err.(validator.ValidationErrors) {
			name := MakeFirstLowerCase(err.StructField())
			switch err.Tag() {
			case "required":
				if name == "nServing" {
					errors[name] = "Invalid target serving"
				} else if name == "recipeID" || name == "recipeId" {
					errors[name] = "Invalid recipe id"
				} else {
					errors[name] = name + " is required"
				}
			case "email":
				errors[name] = name + " should be a valid email"
			case "min":
				if name == "nServing" {
					errors[name] = "Target serving minimum " + err.Param()
				} else {
					errors[name] = name + " minimum " + err.Param() + " characters"
				}

			case "max":
				errors[name] = name + " maximum " + err.Param() + " characters"
			default:
				errors[name] = name + " is invalid"
			}
		}
		return false, errors
	}
	return true, nil

}

//* default validation
func DefaultValidator(c *gin.Context, dataSet interface{}) (bool, interface{}) {
	if err := c.ShouldBind(dataSet); err != nil {
		if _, ok := err.(validator.ValidationErrors); !ok {
			return false, err.Error()
		}

		errors := make(gin.H)

		for _, err := range err.(validator.ValidationErrors) {
			name := MakeFirstLowerCase(err.StructField())
			switch err.Tag() {
			case "required":
				errors[name] = name + " is required"
			case "email":
				errors[name] = name + " should be a valid email"
			case "min":
				errors[name] = name + " minimum " + err.Param() + " characters"
			case "max":
				errors[name] = name + " allowed maximum " + err.Param() + " characters"
			default:
				errors[name] = name + " is invalid"
			}
		}
		return false, errors
	}

	return true, nil
}

func MakeFirstLowerCase(s string) string {
	if len(s) == 0 {
		return s
	}

	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	fmt.Print(string(r))
	return string(r)
}
