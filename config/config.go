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
	debug := loadDefault("DEBUG", "0") != "0"
	debugUsers := loadDefault("DEBUG_USERS", "debug_users.json")
	port := loadDefault("PORT", "8080")
	dbName := loadDefault("DB_NAME", "chrono.db")
	botName := loadStrict("BOT_NAME")
	botEmail := loadStrict("BOT_EMAIL")
	botPassword := loadStrict("BOT_PASSWORD")

	config = &Config{
		Debug:       debug,
		DebugUsers:  debugUsers,
		DbName:      dbName,
		Port:        port,
		BotName:     botName,
		BotEmail:    botEmail,
		BotPassword: botPassword,
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
