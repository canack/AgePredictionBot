package telegram

import (
	"github.com/canack/AgePredictionBot/telegram/handler"
	tele "gopkg.in/telebot.v3"
	_ "image/jpeg"
	_ "image/png"
	"time"
)

var bot *tele.Bot

func setupBotHandlers() {
	bot.Handle("/start", handler.Welcome)
	bot.Handle(tele.OnText, handler.ProcessText)
	bot.Handle("/predict", handler.Predict)
	bot.Handle(tele.OnPhoto, handler.ProcessImage)
	bot.Handle(tele.OnDocument, handler.ProcessImage)
}

func SetupTelegramBot(token string) error {
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	var err error
	bot, err = tele.NewBot(pref)
	if err != nil {
		return err
	}

	setupBotHandlers()

	return nil
}

func StartTelegramBot() {
	bot.Start()
}
