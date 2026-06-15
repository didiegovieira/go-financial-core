package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultEnvironment     = "development"
	defaultHTTPHost        = "0.0.0.0"
	defaultHTTPPort        = 8080
	defaultShutdownTimeout = 10 * time.Second
)

type Settings struct {
	App      AppSettings
	HTTP     HTTPSettings
	Shutdown ShutdownSettings
}

type AppSettings struct {
	Name        string
	Environment string
}

type HTTPSettings struct {
	Host string
	Port int
}

func (s HTTPSettings) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

type ShutdownSettings struct {
	Timeout time.Duration
}

func Load() (Settings, error) {
	httpPort, err := envInt("HTTP_PORT", defaultHTTPPort)
	if err != nil {
		return Settings{}, err
	}

	shutdownTimeout, err := envDuration("SHUTDOWN_TIMEOUT", defaultShutdownTimeout)
	if err != nil {
		return Settings{}, err
	}

	return Settings{
		App: AppSettings{
			Name:        envString("APP_NAME", "payments"),
			Environment: envString("APP_ENV", defaultEnvironment),
		},
		HTTP: HTTPSettings{
			Host: envString("HTTP_HOST", defaultHTTPHost),
			Port: httpPort,
		},
		Shutdown: ShutdownSettings{
			Timeout: shutdownTimeout,
		},
	}, nil
}

func envString(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	return value
}

func envInt(key string, fallback int) (int, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid %s value %q: %w", key, value, err)
	}

	return parsed, nil
}

func envDuration(key string, fallback time.Duration) (time.Duration, error) {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback, nil
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("invalid %s value %q: %w", key, value, err)
	}

	return parsed, nil
}
