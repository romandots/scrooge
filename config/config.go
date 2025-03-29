package config

import "scrooge/utils"

type DatabaseConfig struct {
	Database string
	Username string
	Password string
	Port     string
	Host     string
}

var (
	Socket           = "localhost:8888"
	TelegramBotToken = ""
	TelegramBotDebug = true
	Database         *DatabaseConfig
)

func init() {
	Socket = utils.GetEnvString("SOCKET", Socket)
	TelegramBotToken = utils.GetEnvString("TELEGRAM_BOT_TOKEN", TelegramBotToken)
	TelegramBotDebug = utils.GetEnvBool("TELEGRAM_BOT_DEBUG", TelegramBotDebug)
	Database = &DatabaseConfig{
		Database: utils.GetEnvString("DB_NAME", "scrooge"),
		Username: utils.GetEnvString("DB_USER", "postgres"),
		Password: utils.GetEnvString("DB_PASSWORD", "postgres"),
		Port:     utils.GetEnvString("DB_PORT", "5432"),
	}
}
