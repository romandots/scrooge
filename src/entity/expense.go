package entity

import (
	"fmt"
	"scrooge/messages"
	"scrooge/utils"
	"time"
)

type Expense struct {
	Amount          int `json:"amount" db:"amount"`
	ConvertedAmount int `json:"converted_amount" db:"converted_amount"`
	Rate            *Rate
	Category        string    `json:"category" db:"subject"`
	Receiver        string    `json:"receiver" db:"receiver"`
	Time            time.Time `json:"time" db:"created_at"`
}

func (e *Expense) ConvertAmount() int {
	if e.Rate == nil {
		return e.Amount
	}
	return int(float64(e.Amount) / e.Rate.Rate)
}

func (e *Expense) RublesAmount() int {
	return e.ConvertedAmount
}

func (e *Expense) CurrencyAmount() int {
	return e.Amount
}

func (e *Expense) Currency() string {
	if e.Rate == nil {
		return "₽"
	}
	return e.Rate.Currency
}

func (e *Expense) CurrencyRate() *float64 {
	if e.Rate == nil {
		return nil
	}
	return &e.Rate.Rate
}

func (e *Expense) ToString() string {
	if e.Rate == nil || e.Rate.Rate == 0 {
		return fmt.Sprintf("%d₽", e.Amount)
	}
	return fmt.Sprintf(messages.Price, utils.FormatNumber(e.RublesAmount()), utils.FormatNumber(e.CurrencyAmount()), e.Rate.Currency)
}
