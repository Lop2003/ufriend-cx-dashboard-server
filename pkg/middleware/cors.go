package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CORSConfig() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173, http://localhost:3001",
		AllowMethods:     "GET,POST,PATCH,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	})
}