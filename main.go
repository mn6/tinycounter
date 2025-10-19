package main

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/mn6/tinycounter/config"
	"github.com/mn6/tinycounter/server"
	"github.com/mn6/tinycounter/store"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Load configuration from environment variables
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("config")
		return
	}

	// Pretty print in development
	if cfg.Env != config.Production {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	log.Info().Msg("Loaded configuration.")

	// Create a new Redis store
	store, err := store.NewRedisStore(cfg.RedisAddr, cfg.RedisPass, cfg.RedisDb, cfg.RedisPrefix)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not create redis store.")
	}
	defer store.Close()

	// Initialize and start the Fiber app
	app := server.New(store, cfg)

	// Port to string for fiber
	addr := fmt.Sprintf(":%d", cfg.Port)
	if err := app.Listen(addr, fiber.ListenConfig{DisableStartupMessage: true}); err != nil {
		log.Fatal().Err(err).Msg("server")
	}
}
