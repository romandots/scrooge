package service

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"regexp"
	"scrooge/entity"
	"scrooge/messages"
	"scrooge/postgres"
	"scrooge/utils"
	"strconv"
	"strings"
	"time"
)

func ParseExpenseMessage(text string) (*entity.Expense, bool) {
	if text == "" {
		return nil, false
	}
	lines := strings.Split(text, "\n")

	if len(lines) < 1 || len(lines) > 4 {
		return nil, false
	}

	expense := &entity.Expense{}

	if len(lines) == 1 {
		re := regexp.MustCompile(`^(\d+)\s+([\p{L}]+)\s*(.*)$`)
		matches := re.FindStringSubmatch(text)
		if len(matches) < 2 {
			return nil, false
		}

		amount, err := strconv.Atoi(matches[1])
		if err != nil {
			utils.Error("Failed to parse amount: %v", err)
			return nil, false
		}

		expense.Amount = amount
		expense.Category = matches[2]
		if len(matches) > 3 {
			expense.Receiver = matches[3]
		}

		expense.Time = time.Now()

		return expense, true
	}

	amount, err := strconv.Atoi(lines[0])
	if err != nil {
		utils.Error("Failed to parse amount: %v", err)
		return nil, false
	}

	expense.Amount = amount
	expense.Category = lines[1]

	if len(lines) > 2 {
		expense.Receiver = lines[2]
	}

	if len(lines) > 3 {
		date, err := time.Parse("2006-01-02 15:04", lines[3])
		if err != nil {
			utils.Error("Failed to parse date: %v", err)
			return nil, false
		}
		expense.Time = date
	} else {
		expense.Time = time.Now()
	}

	return expense, true
}

func HandleExpenseMessage(bot *tgbotapi.BotAPI, reply *tgbotapi.MessageConfig, expense *entity.Expense) {
	utils.Debug("Trying to save Expense: %v", expense)
	err := postgres.CreateExpense(expense)
	if err != nil {
		reply.Text = fmt.Sprintf(messages.FailedToSaveExpense, err.Error())
		bot.Send(reply)
		utils.Error(reply.Text)
		return
	}

	utils.Debug("Saved expense: %v", expense)
	reply.Text = fmt.Sprintf(messages.ExpenseSaved, expense.Amount, expense.Category)
	reply.Text += "\n\n"
	todayTotalExpenses, weekCategoryExpenses, monthCategoryExpenses, err := getQuickStatsByCategory(expense.Category)
	if err != nil {
		errorMsg := fmt.Sprintf(messages.FailedToGetQuickStats, err.Error())
		reply.Text += errorMsg
		utils.Error(errorMsg)
	} else {
		reply.Text += fmt.Sprintf(messages.QuickStatsByCategory, utils.FormatDateRussian(time.Now()), todayTotalExpenses, expense.Category, weekCategoryExpenses, monthCategoryExpenses)
	}

	bot.Send(reply)
}

func HandleDelCommand(bot *tgbotapi.BotAPI, reply *tgbotapi.MessageConfig) {
	err := postgres.DeleteLastExpense()
	if err != nil {
		reply.Text = fmt.Sprintf(messages.Error, err.Error())
		utils.Error(reply.Text)
		return
	}
	HandleBalanceCommand(bot, reply, messages.LastExpenseDeleted+"\n\n")
}
