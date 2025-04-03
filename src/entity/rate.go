package entity

type Rate struct {
	Currency string  `json:"currency" db:"currency"`
	Rate     float64 `json:"rate" db:"rate"`
}
