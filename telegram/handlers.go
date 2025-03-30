package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"scrooge/messages"
	"scrooge/service"
	"scrooge/utils"
	"strconv"
	"strings"
	"time"
)

func handleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	if len(message.Text) == 0 {
		return
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "")

	switch message.Text {
	case "/start":
		handleStart(bot, &msg)
	case "/balance":
		handleBalance(bot, &msg, "")
	case "/del":
		handleDelete(bot, &msg)
	default:
		handleCommit(bot, message, &msg)
	}
}

func handleCommit(bot *tgbotapi.BotAPI, message *tgbotapi.Message, reply *tgbotapi.MessageConfig) {
	amount, category, receiver, date, err := parseMessage(message.Text)
	if err != nil {
		reply.Text = messages.FailedToParseMessage
		bot.Send(reply)
		utils.Error("Failed to parse message: %v", err)
		return
	}

	expense, err := service.RecordExpense(amount, category, receiver, date)
	if err != nil {
		reply.Text = fmt.Sprintf(messages.FailedToSaveExpense, err.Error())
		bot.Send(reply)
		utils.Error(reply.Text)
		return
	}

	utils.Debug("Saved expense: %v", expense)
	reply.Text = fmt.Sprintf(messages.ExpenseSaved, expense.Amount, expense.Category)
	reply.Text += "\n\n"
	todayTotalExpenses, weekCategoryExpenses, monthCategoryExpenses, err := service.GetQuickStatsByCategory(expense.Category)
	if err != nil {
		errorMsg := fmt.Sprintf(messages.FailedToGetQuickStats, err.Error())
		reply.Text += errorMsg
		utils.Error(errorMsg)
	} else {
		reply.Text += fmt.Sprintf(messages.QuickStatsByCategory, utils.FormatDateRussian(time.Now()), todayTotalExpenses, expense.Category, weekCategoryExpenses, monthCategoryExpenses)
	}

	bot.Send(reply)
}

func handleStart(bot *tgbotapi.BotAPI, reply *tgbotapi.MessageConfig) {
	reply.Text = messages.StartMessage
	bot.Send(reply)
}

func handleBalance(bot *tgbotapi.BotAPI, reply *tgbotapi.MessageConfig, message string) {
	todayTotalExpenses, weekCategoryExpenses, monthCategoryExpenses, err := service.GetQuickStats()
	if err != nil {
		message += fmt.Sprintf(messages.FailedToGetQuickStats, err.Error())
		utils.Error(message)
	} else {
		message += fmt.Sprintf(messages.QuickStats, utils.FormatDateRussian(time.Now()), todayTotalExpenses, weekCategoryExpenses, monthCategoryExpenses)
	}

	reply.Text = message
	bot.Send(reply)
}

func handleDelete(bot *tgbotapi.BotAPI, reply *tgbotapi.MessageConfig) {
	err := service.DeleteLastExpense()
	if err != nil {
		reply.Text = fmt.Sprintf(messages.Error, err.Error())
		utils.Error(reply.Text)
		return
	}
	handleBalance(bot, reply, messages.LastExpenseDeleted+"\n\n")
}

func parseMessage(text string) (amount int, category string, receiver string, t time.Time, err error) {
	if len(text) == 0 {
		err = fmt.Errorf("empty message")
		return
	}

	lines := strings.Split(text, "\n")
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
