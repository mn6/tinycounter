package routes

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mn6/tinycounter/config"
	"github.com/mn6/tinycounter/handlers/counter"
	"github.com/mn6/tinycounter/handlers/health"
	"github.com/mn6/tinycounter/store"

	"github.com/rs/zerolog/log"
)

// Register all application routes
func Register(app *fiber.App, s store.Store, cfg *config.Config) {
	// Health check endpoint
	app.Get("/health", health.Health)

	log.Info().Msgf("IP Cooldown Duration: %s", cfg.IpCooldownDuration)
	log.Info().Msgf("Image Cache Duration: %s", cfg.ImageCacheDuration)

	// Counter endpoint
	app.Get("/:style/:key", counter.Counter(s, counter.StyleSet, cfg.IpCooldownDuration, cfg.ImageCacheDuration))
}
