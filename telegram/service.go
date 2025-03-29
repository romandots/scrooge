package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
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

func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	bot.Send(msg)
}
