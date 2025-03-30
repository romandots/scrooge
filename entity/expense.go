package entity

import "time"

type Expense struct {
	Amount   int       `json:"amount" db:"amount"`
	Category string    `json:"category" db:"subject"`
	Receiver string    `json:"receiver" db:"receiver"`
	Time     time.Time `json:"time" db:"created_at"`
}
