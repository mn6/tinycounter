package counter

import (
	"os"
	"slices"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/mn6/tinycounter/store"
	"go.yaml.in/yaml/v3"

	"github.com/rs/zerolog/log"
)

type userConfig struct {
	Whitelist []string `yaml:"whitelist"`
	Blacklist []string `yaml:"blacklist"`
}

var users userConfig

func init() {
	// Load users.yaml into users variable
	var err error
	users, err = loadUsersConfig()

	if err != nil {
		users = userConfig{
			Whitelist: []string{},
			Blacklist: []string{},
		}
	}
}

func loadUsersConfig() (userConfig, error) {
	var users userConfig

	data, err := os.ReadFile("configs/users.yaml")
	if err != nil {
		return users, err
	}
	err = yaml.Unmarshal(data, &users)
	if err != nil {
		return users, err
	}
	return users, nil
}

func Counter(s store.Store, styles Styles, cooldownTime time.Duration, cacheTime time.Duration) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Log all headers
		for key, value := range c.GetReqHeaders() {
			log.Debug().Msgf("Header: %s=%s", key, value)
		}

		style := c.Params("style")
		key := c.Params("key")

		// Validate style and key parameters
		isStyleValid := ValidateStyle(style)
		if !isStyleValid {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid style parameter.")
		}
		if key == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Key parameter is required.")
		}
		if len(key) > 64 {
			key = key[:64]
		}

		if len(users.Whitelist) > 0 {
			// If there is a non-empty whitelist, enforce it
			allowed := slices.Contains(users.Whitelist, key)
			if !allowed {
				return c.Status(fiber.StatusForbidden).SendString("Access denied: key not in whitelist.")
			}
		} else if len(users.Blacklist) > 0 {
			// If there is no whitelist but there is a blacklist, enforce it
			denied := slices.Contains(users.Blacklist, key)
			if denied {
				return c.Status(fiber.StatusForbidden).SendString("Access denied: key is blacklisted.")
			}
		}

		// Get the current counter value
		counterValue, err := s.GetCounter(c.Context(), key)
		if err != nil {
			counterValue = 1
		}

		// Capture request-scoped values before launching the goroutine
		ip := c.IP()
		ctx := c.Context()

		// Attempt to set cooldown for this IP and key in a separate goroutine
		go func() {
			canIncrement, err := s.TrySetCooldown(ctx, key, ip, cooldownTime)
			if canIncrement && err == nil {
				_, _ = s.IncrementCounter(ctx, key)
			}
		}()

		// Generate the counter image
		imgData := GenerateCounter(s, key, styles, style, counterValue, cacheTime)

		// Set appropriate headers and send the image
		c.Set("Content-Type", "image/png")
		return c.Status(fiber.StatusOK).Send(imgData)
	}
}
