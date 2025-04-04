package messages

var (
	Help                  = "Используйте следующий формат (через пробел):\n*сумма категория* (одно слово) *получатель* (опционально)\n\nИли в несколько строк:\n*сумма*\n*категория* _(можно несколько слов)_\n*получатель* _(опционально)_\n*дата и время* _(опционально)_\n\nНапример:\n*1000 продукты магазин у дома*\n\nИли:\n*1000*\n*хоз товары*\n*магазин у дома*\n*2021-01-01 12:00*\n\nТакже вы можете указать валюту, предварительно указав ее курс:\n*100usd продукты магазин*\n\nЧтобы указать курс, отправьте сообщение следующего формата:\n*курс usd 0.012*\n_(где 0.012 - это курс рубля к доллару)_"
	StartCommand          = "Что это за бот?"
	BalanceCommand        = "Баланс"
	DelCommand            = "Удаление последней записи"
	RatesCommand          = "Посмотреть установленные курсы валют"
	StartMessage          = "Привет. Я помогу тебе вести учет расходов.\n\n" + Help
	Error                 = "Возникла ошибка.\n%s"
	ExpenseSaved          = "Сохранили расход: %s на %s"
	ExpenseSavedIn        = "Сохранили расход: %s на %s в %s"
	FailedToGetQuickStats = "Не удалось получить статистику: %s"
	FailedToParseMessage  = "Сообщение не распознано.\n\n" + Help
	FailedToSaveExpense   = "Не удалось сохранить запись: %s"
	QuickStatsByCategory  = "Расходы за сегодня, %[1]s: %[2]s₽\nРасходы на %[3]s за неделю: %[4]s₽\nРасходы на %[3]s за месяц: %[5]s₽"
	QuickStats            = "Расходы за сегодня, %[1]s: %[2]s₽\nРасходы за неделю: %[3]s₽\nРасходы за месяц: %[4]s₽"
	LastExpenseDeleted    = "Последний расход удален"
	RateSet               = "Курс установлен: %s"
	FailedToSaveCache     = "Не удалось сохранить курс в кэше: %s"
	FailedToGetRates      = "Не удалось получить курсы валют: %s"
	Rates                 = "Текущие курсы валют:\n"
	RateLine              = "1 %[1]s = %.3[3]f₽ / 1₽ = %.3[2]f %[1]s\n"
	Price                 = "%s₽ (%s %s)"
)
