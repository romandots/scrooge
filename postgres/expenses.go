package postgres

import (
	"scrooge/entity"
	"scrooge/utils"
	"strings"
	"time"
)

type FilterExpenses struct {
	Subject  *string
	Receiver *string
	FromDate *time.Time
	ToDate   *time.Time
}

func GetTotalExpensesToday() (int, error) {
	startOfToday := utils.StartOfDay(time.Now())
	return getExpensesTotal(&FilterExpenses{
		FromDate: &startOfToday,
	})
}

func GetWeekExpensesBySubject(subject *string) (int, error) {
	startOfThisWeek := utils.StartOfWeek(time.Now())
	return getExpensesTotal(&FilterExpenses{
		Subject:  subject,
		FromDate: &startOfThisWeek,
	})
}

func GetMonthExpensesBySubject(subject *string) (int, error) {
	startOfThisMonth := utils.StartOfMonth(time.Now())
	return getExpensesTotal(&FilterExpenses{
		Subject:  subject,
		FromDate: &startOfThisMonth,
	})
}

func CreateExpense(expense *entity.Expense) error {
	sql, args, err := Sq.
		Insert("expenses").
		SetMap(map[string]interface{}{
			"amount":     expense.Amount,
			"subject":    expense.Category,
			"receiver":   expense.Receiver,
			"created_at": expense.Time,
		}).
		ToSql()
	if err != nil {
		return err
	}

	return exec(sql, args...)
}

func DeleteLastExpense() error {
	sql, args, err := Sq.
		Delete("expenses").
		Where("created_at = (SELECT MAX(created_at) FROM expenses)").
		ToSql()
	if err != nil {
		return err
	}

	return exec(sql, args...)
}

func getExpenses(filter FilterExpenses) ([]entity.Expense, error) {
	builder := Sq.
		Select("amount", "subject", "receiver", "created_at").
		From("expenses").
		Where("true")

	if filter.Subject != nil {
		builder = builder.Where("lower(subject) = ?", strings.ToLower(*filter.Subject))
	}

	if filter.Receiver != nil {
		builder = builder.Where("lower(receiver) = ?", strings.ToLower(*filter.Receiver))
	}

	if filter.ToDate != nil {
		builder = builder.Where("created_at <= ?", *filter.ToDate)
	}

	sql, args, err := builder.ToSql()
	rows, err := query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []entity.Expense
	for rows.Next() {
		var expense entity.Expense
		err := rows.Scan(&expense.Amount, &expense.Category, &expense.Receiver, &expense.Time)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}

	return expenses, nil
}

func getExpensesTotal(filter *FilterExpenses) (int, error) {
	builder := Sq.
		Select("COALESCE(SUM(amount), 0)").
		From("expenses").
		Where("true")

	if filter.Subject != nil {
		builder = builder.Where("lower(subject) = ?", strings.ToLower(*filter.Subject))
	}

	if filter.Receiver != nil {
		builder = builder.Where("lower(receiver) = ?", strings.ToLower(*filter.Receiver))
	}

	if filter.FromDate != nil {
		builder = builder.Where("created_at >= ?", *filter.FromDate)
	}

	if filter.ToDate != nil {
		builder = builder.Where("created_at <= ?", *filter.ToDate)
	}

	sql, args, err := builder.ToSql()
	rows, err := query(sql, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var total int
	for rows.Next() {
		err := rows.Scan(&total)
		if err != nil {
			return 0, err
		}
	}

	return total, nil
}
