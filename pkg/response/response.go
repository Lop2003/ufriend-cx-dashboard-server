package response

import "github.com/gofiber/fiber/v2"

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func OK(c *fiber.Ctx, message string, data interface{}) error {
	return c.JSON(APIResponse{Success: true, Message: message, Data: data})
}

func Error(c *fiber.Ctx, statusCode int, message string, errCode string) error {
	return c.Status(statusCode).JSON(APIResponse{Success: false, Message: message, Error: errCode})
}

func ErrorServer(c *fiber.Ctx, message string, err error) error {
	return c.Status(fiber.StatusInternalServerError).JSON(APIResponse{Success: false, Message: message, Error: err.Error()})
}