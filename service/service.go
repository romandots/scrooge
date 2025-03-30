package service

import (
	"scrooge/entity"
	"scrooge/postgres"
	"scrooge/utils"
	"time"
)

func DeleteLastExpense() error {
	return postgres.DeleteLastExpense()
}

func RecordExpense(amount int, category string, receiver string, date time.Time) (expense *entity.Expense, err error) {
	expense = &entity.Expense{
		Amount:   amount,
		Category: category,
		Receiver: receiver,
		Time:     date,
	}

	utils.Debug("Trying to save Expense: %v", expense)
	err = postgres.CreateExpense(expense)

	return
}

func GetQuickStats() (todayTotalExpenses, weekTotalExpenses, monthTotalExpenses int, err error) {
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

func GetQuickStatsByCategory(category string) (todayTotalExpenses, weekCategoryExpenses, monthCategoryExpenses int, err error) {
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
