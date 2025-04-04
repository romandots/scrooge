package utils

import (
	"github.com/dustin/go-humanize"
	"strings"
)

func FormatNumber(number int) string {
	formattedNumber := humanize.Comma(int64(number))
	return strings.ReplaceAll(formattedNumber, ",", " ")
}
