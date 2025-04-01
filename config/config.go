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

func init() {
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
}
