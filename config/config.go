package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Environment represents the environment the application is running in
type Environment string

// Possible Environment values
const (
	Development Environment = "dev"
	Production  Environment = "prod"
)

// Environment variable keys for consistency
const (
	EnvKeyPort               = "APP_PORT"
	EnvKeyEnv                = "ENV"
	EnvKeyIpCooldownDuration = "IP_COOLDOWN_DURATION"
	EnvKeyImageCacheDuration = "IMAGE_CACHE_DURATION"

	EnvKeyRedisAddr   = "REDIS_ADDR_IN_CONTAINER"
	EnvKeyRedisPass   = "REDIS_PASSWORD"
	EnvKeyRedisDb     = "REDIS_DB"
	EnvKeyRedisPrefix = "REDIS_PREFIX"
)

// Config holds the application configuration
type Config struct {
	Port               int
	Env                Environment
	IpCooldownDuration time.Duration
	ImageCacheDuration time.Duration

	RedisAddr   string
	RedisPass   string
	RedisDb     int
	RedisPrefix string
}

// DefaultConfig provides default configuration values
var DefaultConfig = &Config{
	Port:               4392,
	Env:                Development,
	IpCooldownDuration: 30 * time.Second,
	ImageCacheDuration: 1 * time.Minute,

	RedisAddr:   "redis:6379",
	RedisPass:   "",
	RedisDb:     0,
	RedisPrefix: "tinycounter",
}

// Load reads configuration from environment variables and returns a Config struct
// It falls back to default values if environment variables are not set
func Load() (*Config, error) {
	Port := getEnvAsInt(EnvKeyPort, DefaultConfig.Port)

	Env := Environment(os.Getenv(EnvKeyEnv))
	if Env == "" {
		Env = DefaultConfig.Env
	}

	RedisAddr := os.Getenv(EnvKeyRedisAddr)
	if RedisAddr == "" {
		RedisAddr = DefaultConfig.RedisAddr
	}

	RedisPass := os.Getenv(EnvKeyRedisPass)
	if RedisPass == "" {
		RedisPass = DefaultConfig.RedisPass
	}

	RedisDb := getEnvAsInt(EnvKeyRedisDb, DefaultConfig.RedisDb)

	RedisPrefix := os.Getenv(EnvKeyRedisPrefix)
	if RedisPrefix == "" {
		RedisPrefix = DefaultConfig.RedisPrefix
	}

	if Port <= 0 {
		return nil, fmt.Errorf("Invalid %s: %d", EnvKeyPort, Port)
	}

	if Env != Development && Env != Production {
		return nil, fmt.Errorf("Invalid %s: %s", EnvKeyEnv, Env)
	}

	IpCooldownDuration := getEnvAsDuration(EnvKeyIpCooldownDuration, DefaultConfig.IpCooldownDuration)
	ImageCacheDuration := getEnvAsDuration(EnvKeyImageCacheDuration, DefaultConfig.ImageCacheDuration)

	if RedisAddr == "" {
		return nil, fmt.Errorf("Invalid %s: cannot be empty", EnvKeyRedisAddr)
	}

	return &Config{
		Port:               Port,
		Env:                Env,
		IpCooldownDuration: IpCooldownDuration,
		ImageCacheDuration: ImageCacheDuration,

		RedisAddr:   RedisAddr,
		RedisPass:   RedisPass,
		RedisDb:     RedisDb,
		RedisPrefix: RedisPrefix,
	}, nil
}

// getEnvAsInt retrieves an environment variable as an integer,
// or returns a default value if not set or invalid
func getEnvAsInt(key string, defaultVal int) int {
	if val, exists := os.LookupEnv(key); exists {
		returnVal, err := strconv.Atoi(val)
		if err != nil {
			return defaultVal
		} else {
			return returnVal
		}
	} else {
		return defaultVal
	}
}

// getEnvAsDuration retrieves an environment variable as a time.Duration,
// or returns a default value if not set or invalid
func getEnvAsDuration(key string, defaultVal time.Duration) time.Duration {
	if val, exists := os.LookupEnv(key); exists {
		returnVal, err := time.ParseDuration(val)
		// log err
		if err != nil {
			return defaultVal
		} else {
			return returnVal
		}
	} else {
		return defaultVal
	}
}
