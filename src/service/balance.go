package service

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"scrooge/messages"
	"scrooge/postgres"
	"scrooge/utils"
	"time"
)

func HandleBalanceCommand(bot *tgbotapi.BotAPI, reply *tgbotapi.MessageConfig, message string) {
	todayTotalExpenses, weekCategoryExpenses, monthCategoryExpenses, err := getQuickStats()
	if err != nil {
		message += fmt.Sprintf(messages.FailedToGetQuickStats, err.Error())
		utils.Error(message)
	} else {
		message += fmt.Sprintf(messages.QuickStats, utils.FormatDateRussian(time.Now()), todayTotalExpenses, weekCategoryExpenses, monthCategoryExpenses)
	}

	reply.Text = message
	bot.Send(reply)
}

func getQuickStats() (todayTotalExpenses, weekTotalExpenses, monthTotalExpenses int, err error) {
	todayTotalExpenses, err = postgres.GetTotalExpensesToday()
	if err != nil {
		return
	}

	weekTotalExpenses, err = postgres.GetWeekExpensesBySubject(nil)
	if err != nil {
		return
	}

	monthTotalExpenses, err = postgres.GetMonthExpensesBySubject(nil)
	if err != nil {
		return
	}

	return
}

func getQuickStatsByCategory(category string) (todayTotalExpenses, weekCategoryExpenses, monthCategoryExpenses int, err error) {
	todayTotalExpenses, err = postgres.GetTotalExpensesToday()
	if err != nil {
		return
	}

	weekCategoryExpenses, err = postgres.GetWeekExpensesBySubject(&category)
	if err != nil {
		return
	}

	monthCategoryExpenses, err = postgres.GetMonthExpensesBySubject(&category)
	if err != nil {
		return
	}

	return
}
