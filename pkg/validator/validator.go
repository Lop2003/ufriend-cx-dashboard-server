package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ValidateStruct(s interface{}) []ValidationError {
	var errs []ValidationError

	err := validate.Struct(s)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			errs = append(errs, ValidationError{
				Field:   toSnakeCase(e.Field()),
				Message: formatMessage(e),
			})
		}
	}

	return errs
}

func ValidateBody(c *fiber.Ctx, out interface{}) error {
	if err := c.BodyParser(out); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid request body",
			"error":   "ERR_PARSE",
		})
	}

	validationErrors := ValidateStruct(out)
	if len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "validation failed",
			"error":   "ERR_VALIDATION",
			"data":    validationErrors,
		})
	}

	return nil
}

func formatMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", toSnakeCase(e.Field()))
	case "gte":
		return fmt.Sprintf("%s must be at least %s", toSnakeCase(e.Field()), e.Param())
	case "lte":
		return fmt.Sprintf("%s must be at most %s", toSnakeCase(e.Field()), e.Param())
	case "min":
		return fmt.Sprintf("%s must be at least %s", toSnakeCase(e.Field()), e.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s", toSnakeCase(e.Field()), e.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", toSnakeCase(e.Field()), e.Param())
	default:
		return fmt.Sprintf("%s is invalid", toSnakeCase(e.Field()))
	}
}

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}