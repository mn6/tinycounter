package server

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/mn6/tinycounter/config"
	"github.com/rs/zerolog/log"
)

func RequestLogger(env *config.Environment) func(fiber.Ctx) error {
	return func(c fiber.Ctx) error {
		// Start timer
		start := time.Now()

		// Limit user agent length to 200 chars
		userAgent := string(c.Get("User-Agent"))
		if len(userAgent) > 200 {
			userAgent = userAgent[:200]
		}

		// Create logger with request fields
		reqLogger := log.With().
			Str("method", c.Method()).
			Str("path", c.Path()).
			Str("remote_addr", c.IP()).
			Str("user_agent", userAgent).
			Logger()

		// Store logger in context locals for handlers to use
		c.Locals("logger", &reqLogger)

		// Proceed or log error
		if err := c.Next(); err != nil {
			reqLogger.Error().
				Err(err).
				Int("status", c.Response().StatusCode()).
				Dur("latency", time.Since(start)).
				Msg("request error")
			return err
		}

		// Log request completion
		reqLogger.Info().
			Int("status", c.Response().StatusCode()).
			Dur("latency", time.Since(start)).
			Msg("request completed")
		return nil
	}
}
