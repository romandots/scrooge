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

func ParseExpenseMessage(text string) (*entity.Expense, bool, error) {
	if text == "" {
		return nil, false, nil
	}
	lines := strings.Split(text, "\n")

	if len(lines) < 1 || len(lines) > 4 {
		return nil, false, nil
	}

	expense := &entity.Expense{
		Rate: &entity.Rate{},
	}
	var (
		amount   int
		category string
		receiver string
		currency string
		rate     float64
		err      error
		date     = time.Now()
	)

	if len(lines) == 1 {
		// Single line commit
		re := regexp.MustCompile(`^(\d+)([\S\w\p{L}]*)\s+([\p{L}]+)\s*(.*)$`)
		matches := re.FindStringSubmatch(text)
		if len(matches) < 2 {
			return nil, false, nil
		}

		amount, err = strconv.Atoi(strings.Trim(matches[1], " "))
		if err != nil {
			return nil, false, utils.Error("Failed to parse amount: %v", err)
		}

		currency = strings.Trim(matches[2], " ")
		category = strings.Trim(matches[3], " ")
		if len(matches) > 4 {
			receiver = strings.Trim(matches[4], " ")
		}

	} else {
		// Multiline commit
		amountRaw := strings.Trim(lines[0], " ")
		currencyRegex := regexp.MustCompile(`^(\d+)\s*([\S\w\p{L}]*)$`)
		matches := currencyRegex.FindStringSubmatch(amountRaw)
		if len(matches) < 2 {
			return nil, false, nil
		}
		amount, err = strconv.Atoi(matches[1])
		if err != nil {
			return nil, false, utils.Error("Failed to parse amount: %v", err)
		}
		currency = strings.Trim(matches[2], " ")
		category = strings.Trim(lines[1], " ")
		if len(lines) > 2 {
			receiver = strings.Trim(lines[2], " ")
		}

		if len(lines) > 3 {
			customDate, err := time.Parse("2006-01-02 15:04", strings.Trim(lines[3], " "))
			if err != nil {
				return nil, false, utils.Error("Failed to parse date: %v", err)
			}
			date = customDate
		}
	}

	currency = strings.ToUpper(currency)
	if currency != "" {
		rate, err = getRate(currency)
		if err != nil || rate == 0 {
			return nil, false, utils.Error("Курс %s не установлен", currency)
		}
	}

	expense = &entity.Expense{
		Amount:          amount,
		ConvertedAmount: amount,
		Rate:            &entity.Rate{},
		Category:        category,
		Receiver:        receiver,
		Time:            date,
	}

	if currency != "" {
		expense.Rate.Rate = rate
		expense.Rate.Currency = currency
		expense.ConvertedAmount = expense.ConvertAmount()
	}

	return expense, true, nil
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
	if expense.Receiver != "" {
		reply.Text = fmt.Sprintf(messages.ExpenseSavedIn, expense.ToString(), expense.Category, expense.Receiver)
	} else {
		reply.Text = fmt.Sprintf(messages.ExpenseSaved, expense.ToString(), expense.Category)
	}
	reply.Text += "\n\n"
	todayTotalExpenses, weekCategoryExpenses, monthCategoryExpenses, err := getQuickStatsByCategory(expense.Category)
	if err != nil {
		errorMsg := fmt.Sprintf(messages.FailedToGetQuickStats, err.Error())
		reply.Text += errorMsg
		utils.Error(errorMsg)
	} else {
		reply.Text += fmt.Sprintf(messages.QuickStatsByCategory, utils.FormatDateRussian(time.Now()), utils.FormatNumber(todayTotalExpenses), expense.Category, utils.FormatNumber(weekCategoryExpenses), utils.FormatNumber(monthCategoryExpenses))
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
