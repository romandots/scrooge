package service

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"regexp"
	"scrooge/cache"
	"scrooge/entity"
	"scrooge/messages"
	"scrooge/utils"
	"strconv"
	"strings"
)

func ParseRateMessage(text string) (*entity.Rate, bool) {
	if text == "" {
		return nil, false
	}

	text = strings.ToUpper(strings.Trim(text, " "))
	rateMatch := regexp.MustCompile(`^КУРС\s+([\d,\.]+)\s*(.*)$`)
	rateIndex := 1
	currencyIndex := 2
	if !rateMatch.MatchString(text) {
		rateMatch = regexp.MustCompile(`^КУРС\s+([\D]+)\s*([\d,\.]+)$`)
		rateIndex = 2
		currencyIndex = 1
	}
	if !rateMatch.MatchString(text) {
		return nil, false
	}

	matches := rateMatch.FindStringSubmatch(text)
	if len(matches) < 3 {
		return nil, false
	}

	float, err := strconv.ParseFloat(strings.ReplaceAll(matches[rateIndex], ",", "."), 64)
	if err != nil {
		utils.Error("Failed to parse rate: %v", err)
		return nil, false
	}

	return &entity.Rate{
		Rate:     float,
		Currency: strings.Trim(matches[currencyIndex], " "),
	}, true
}

func HandleRateMessage(bot *tgbotapi.BotAPI, reply *tgbotapi.MessageConfig, rate *entity.Rate) {
	cacheKey := cacheKey(rate.Currency)
	err := cache.Set(cacheKey, rate.Rate)
	if err != nil {
		reply.Text = fmt.Sprintf(messages.FailedToSaveCache, err.Error())
		bot.Send(reply)
		utils.Error(reply.Text)
		return
	}
	reverseRate := 1 / rate.Rate
	reply.Text = fmt.Sprintf(messages.RateSet, fmt.Sprintf(messages.RateLine, rate.Currency, rate.Rate, reverseRate))
	bot.Send(reply)
}

func HandleRatesCommand(bot *tgbotapi.BotAPI, reply *tgbotapi.MessageConfig) {
	ratesRaw, err := cache.GetAll("rate:*")
	if err != nil {
		reply.Text = fmt.Sprintf(messages.FailedToGetRates, err.Error())
		utils.Error(reply.Text)
		return
	}

	message := messages.Rates
	for cacheKey, rateRaw := range ratesRaw {
		currency := getCurrencyFromCacheKey(cacheKey)
		rate, err := strconv.ParseFloat(rateRaw, 64)
		if err != nil {
			utils.Error("Failed to parse rate: %v", err)
			continue
		}
		reverseRate := 1 / rate
		message += fmt.Sprintf(messages.RateLine, currency, rate, reverseRate)
	}

	reply.Text = message
	bot.Send(reply)
}

func getRate(currency string) (float64, error) {
	cacheKey := cacheKey(currency)
	rateRaw, err := cache.Get(cacheKey)
	if rateRaw == "" {
		return 0, nil
	}

	rate, err := strconv.ParseFloat(rateRaw, 64)
	if err != nil {
		return 0, err
	}

	return rate, nil
}

func cacheKey(currency string) string {
	return fmt.Sprintf("rate:%s", currency)
}

func getCurrencyFromCacheKey(cacheKey string) string {
	return strings.TrimPrefix(cacheKey, "rate:")
}
