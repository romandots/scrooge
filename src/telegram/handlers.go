package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"scrooge/messages"
	"scrooge/service"
)

func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	if len(message.Text) == 0 {
		return
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "")
	msg.ParseMode = "Markdown"

	switch message.Text {
	case "/start":
		handleStartCommand(bot, &msg)
		return
	case "/balance":
		service.HandleBalanceCommand(bot, &msg, "")
		return
	case "/del":
		service.HandleDelCommand(bot, &msg)
		return
	case "/rates":
		service.HandleRatesCommand(bot, &msg)
		return
	}

	rate, found := service.ParseRateMessage(message.Text)
	if found {
		service.HandleRateMessage(bot, &msg, rate)
		return
	}

	expense, found, err := service.ParseExpenseMessage(message.Text)
	if err != nil {
		msg.Text = err.Error()
		bot.Send(msg)
		return
	}

	if found {
		service.HandleExpenseMessage(bot, &msg, expense)
		return
	}

	msg.Text = messages.FailedToParseMessage
	bot.Send(msg)
}

func handleStartCommand(bot *tgbotapi.BotAPI, reply *tgbotapi.MessageConfig) {
	reply.Text = messages.StartMessage
	bot.Send(reply)
}
