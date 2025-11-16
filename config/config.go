package config

import (
	"log"
	"log/slog"
	"os"
)

type Config struct {
	Debug       bool
	DebugUsers  string
	Port        string
	DbName      string
	BotPassword string
	BotEmail    string
	BotName     string
	Banner      string
	SentryUrl   string
	AworkApiKey string
}

var config *Config

func GetConfig() *Config {
	return config
}

func NewConfigFromEnv() *Config {
	config = &Config{
		Debug:       loadDefault("DEBUG", "0") == "1",
		DebugUsers:  loadDefault("DEBUG_USERS", "debug_users.json"),
		DbName:      loadDefault("DB_NAME", "chrono.db"),
		Port:        loadDefault("PORT", "8080"),
		BotName:     loadStrict("BOT_NAME"),
		BotEmail:    loadStrict("BOT_EMAIL"),
		BotPassword: loadStrict("BOT_PASSWORD"),
		SentryUrl: loadIf(
			"SENTRY_URL",
			func() bool { return loadDefault("DEBUG", "0") == "0" },
		),
		AworkApiKey: loadDefault("AWORK_API_KEY", ""),
	}

	slog.Info("Config loaded")

	return config
}

func loadIf(val string, fn func() bool) string {
	if fn() {
		return loadDefault(val, "")
	}

	return ""
}

func loadStrict(envVar string) string {
	valEnv := os.Getenv(envVar)
	if valEnv == "" {
		log.Fatalf("Environment variable \"%v\" is missing", envVar)
	}

	return valEnv
}

func loadDefault(envVar, defaultVal string) string {
	valEnv := os.Getenv(envVar)
	if valEnv == "" {
		valEnv = defaultVal
	}

	return valEnv
}
