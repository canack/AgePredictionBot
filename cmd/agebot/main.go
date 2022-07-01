package main

import (
	"github.com/canack/AgePredictionBot/recognize"
	"github.com/canack/AgePredictionBot/telegram"
	"log"
	"os"
)

var token string

func main() {

	if tokenEnv := os.Getenv("BOT_TOKEN"); tokenEnv == "" {
		panic("Token is not declared.\nPlease attach your token as environment variable. Eg: BOT_TOKEN='token'")
	} else {
		token = tokenEnv
	}

	log.Println("Bot started")

	startBot()
}

func startBot() {
	if err := telegram.SetupTelegramBot(token); err != nil {
		panic(err)
	}

	if err := recognize.SetupRekognition(); err != nil {
		panic(err)
	}

	telegram.StartTelegramBot()
}
