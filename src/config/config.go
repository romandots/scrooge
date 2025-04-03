package config

import "scrooge/utils"

type DatabaseConfig struct {
	Database string
	Username string
	Password string
	Port     string
	Host     string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

var (
	Socket           = "localhost:8888"
	TelegramBotToken = ""
	TelegramBotDebug = false
	Database         *DatabaseConfig
	Redis            *RedisConfig
)

func InitConfig() {
	Socket = utils.GetEnvString("SOCKET", Socket)
	TelegramBotToken = utils.GetEnvString("TELEGRAM_BOT_TOKEN", TelegramBotToken)
	TelegramBotDebug = utils.GetEnvBool("TELEGRAM_BOT_DEBUG", TelegramBotDebug)
	Database = &DatabaseConfig{
		Host:     utils.GetEnvString("DB_HOST", "localhost"),
		Port:     utils.GetEnvString("DB_PORT", "5432"),
		Database: utils.GetEnvString("DB_NAME", "scrooge"),
		Username: utils.GetEnvString("DB_USER", "postgres"),
		Password: utils.GetEnvString("DB_PASSWORD", "postgres"),
	}
	Redis = &RedisConfig{
		Host:     utils.GetEnvString("REDIS_HOST", "localhost"),
		Port:     utils.GetEnvString("REDIS_PORT", "6379"),
		Password: utils.GetEnvString("REDIS_PASSWORD", ""),
		DB:       utils.GetEnvInt("REDIS_DB", 0),
	}

	if TelegramBotToken == "" {
		utils.Error("Telegram bot token is not set")
		panic("Telegram bot token is not set")
	}

	utils.Debug("Socket: %s", Socket)
	utils.Debug("Telegram Bot Token: %s", TelegramBotToken)
	utils.Debug("Telegram Bot Debug: %v", TelegramBotDebug)
	utils.Debug("Database Host: %s", Database.Host)
	utils.Debug("Database Port: %s", Database.Port)
	utils.Debug("Database Name: %s", Database.Database)
	utils.Debug("Database Username: %s", Database.Username)
	utils.Debug("Database Password: %s", Database.Password)
	utils.Debug("Redis Host: %s", Redis.Host)
	utils.Debug("Redis Port: %s", Redis.Port)
}
