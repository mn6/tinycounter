package health

import "github.com/gofiber/fiber/v3"

// Health returns a small JSON payload for health checks
func Health(c fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"service": "tinycounter",
	})
}
