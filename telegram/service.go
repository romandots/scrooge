package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"scrooge/config"
	"scrooge/utils"
)

func StartBot() {
	utils.Info("Starting Telegram bot...")

	bot, err := tgbotapi.NewBotAPI(config.TelegramBotToken)
	if err != nil {
		utils.Error("Failed to start telegram bot: " + err.Error())
		panic(err)
	}
	bot.Debug = config.TelegramBotDebug

	// register commands
	commands := []tgbotapi.BotCommand{
		{Command: "start", Description: "Start the bot"},
		{Command: "balance", Description: "Get your balance"},
		{Command: "del", Description: "Delete the last expense"},
	}
	_, err = bot.Request(tgbotapi.NewSetMyCommands(commands...))
	if err != nil {
		utils.Error("Failed to set bot commands: " + err.Error())
		panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	if err != nil {
		utils.Error("Failed to receive updates from Telegram: " + err.Error())
		panic(err)
	}

	for update := range updates {
		if update.Message != nil {
			utils.Debug("Telegram message received: %v", update.Message)
			handleMessage(bot, update.Message)
		}
	}
}
