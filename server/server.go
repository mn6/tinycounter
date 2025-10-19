package server

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	"github.com/mn6/tinycounter/config"
	"github.com/mn6/tinycounter/routes"
	"github.com/mn6/tinycounter/store"
)

// New creates and configures a Fiber app with routes registered.
func New(s store.Store, cfg *config.Config) *fiber.App {
	// Establish new Fiber with sonic encoder/decoder for better performance
	// May not use json but good to have available
	app := fiber.New(fiber.Config{
		AppName:     "TinyCounter",
		JSONEncoder: sonic.Marshal,
		JSONDecoder: sonic.Unmarshal,
	})

	// -- Register middleware --
	// Request logging, passing in the environment for logging pretty print
	app.Use(RequestLogger(&cfg.Env))

	// -- Register routes --
	routes.Register(app, s, cfg)
	return app
}
