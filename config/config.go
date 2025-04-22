package config

import (
	"log"
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
}

var config *Config

func GetConfig() *Config {
	return config
}

func NewConfigFromEnv() *Config {
	config = &Config{
		Debug:       loadDefault("DEBUG", "0") != "0",
		DebugUsers:  loadDefault("DEBUG_USERS", "debug_users.json"),
		DbName:      loadDefault("DB_NAME", "chrono.db"),
		Port:        loadDefault("PORT", "8080"),
		BotName:     loadStrict("BOT_NAME"),
		BotEmail:    loadStrict("BOT_EMAIL"),
		BotPassword: loadStrict("BOT_PASSWORD"),
	}

	log.Println("Config loaded")

	return config
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
