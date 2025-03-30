package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"scrooge/config"
	"scrooge/messages"
	"scrooge/service"
	"scrooge/utils"
	"strconv"
	"strings"
	"time"
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
	if len(message.Text) == 0 {
		return
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "")
	amount, category, receiver, date, err := parseMessage(message.Text)
	if err != nil {
		msg.Text = messages.FailedToParseMessage
		bot.Send(msg)
		utils.Error("Failed to parse message: %v", err)
		return
	}

	expense, err := service.RecordExpense(amount, category, receiver, date)
	if err != nil {
		msg.Text = fmt.Sprintf(messages.FailedToSaveExpense, err.Error())
		bot.Send(msg)
		utils.Error(msg.Text)
		return
	}

	utils.Debug("Saved expense: %v", expense)
	msg.Text = fmt.Sprintf(messages.ExpenseSaved, expense.Amount, expense.Category)
	msg.Text += "\n\n"
	todayTotalExpenses, weekCategoryExpenses, monthCategoryExpenses, err := service.GetQuickStats(expense.Category)
	if err != nil {
		errorMsg := fmt.Sprintf(messages.FailedToGetQuickStats, err.Error())
		msg.Text += errorMsg
		utils.Error(errorMsg)
	} else {
		msg.Text += fmt.Sprintf(messages.QuickStats, utils.FormatDateRussian(time.Now()), todayTotalExpenses, expense.Category, weekCategoryExpenses, monthCategoryExpenses)
	}

	bot.Send(msg)
}

func parseMessage(text string) (amount int, category string, receiver string, t time.Time, err error) {
	if len(text) == 0 {
		err = fmt.Errorf("empty message")
		return
	}

	lines := strings.Split(text, "\n")
	utils.Debug("Message has %d lines", len(lines))
	if len(lines) < 2 || len(lines) > 4 {
		err = fmt.Errorf("wrong format")
		return
	}

	amount, err = strconv.Atoi(lines[0])
	if err != nil {
		return
	}

	category = lines[1]

	if len(lines) > 2 {
		receiver = lines[2]
	}

	if len(lines) > 3 {
		t, err = time.Parse("2006-01-02 15:04", lines[3])
		if err != nil {
			return
		}
	} else {
		t = time.Now()
	}

	return
}
