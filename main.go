package main

import (
	"log"
	"os"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"

	"fsmExample/callbacks"
	"fsmExample/consts"
	"fsmExample/database"
)

// all "Messages" handler
func MainMessageProcessor(api tgbotapi.BotAPI, update tgbotapi.Update) {
	state := database.GetCurrentState(update)

	switch {

	// handle commands
	case state == database.InitialState: // empty state not initialized yet
		if update.Message.Text == "/start" {
			go callbacks.CmdStart(api, update)
		} else if update.Message.Text == "/stats" {
			go callbacks.CmdStats(api, update)
		}

	// handle other text updates
	case state == consts.Name:
		go callbacks.ProcessName(api, update)
	case state == consts.Gender:
		go callbacks.ProcessGender(api, update)
	}
}

// callback queries handler
func MainCallbackQueryProcessor(api tgbotapi.BotAPI, update tgbotapi.Update) {
	state := database.GetCurrentState(update)

	switch {
	case state == consts.ItemChoose:
		go callbacks.ClickOnItem(api, update)

	// check if Data has prefix ukey:
	case state == database.InitialState && strings.HasPrefix(
		update.CallbackQuery.Data,
		consts.StatsCallbackDataPrefix+":",
	):
		go callbacks.ProcessAbout(api, update)
	}
}

// updates entry point
func HandleUpdateBase(api tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := api.GetUpdatesChan(u)

	if err != nil {
		panic(err)
	}

	for update := range updates {
		switch {
		case update.Message != nil:
			go MainMessageProcessor(api, update)
		case update.CallbackQuery != nil:
			go MainCallbackQueryProcessor(api, update)
		}
	}
}

func main() {
	// entry and initialization point
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	api, err := tgbotapi.NewBotAPI(os.Getenv("GO_BOT_TOKEN"))

	if err != nil {
		log.Panic(err)
	}

	api.Debug = true

	log.Printf("Running @%s[%d]", api.Self.UserName, api.Self.ID)

	HandleUpdateBase(*api) // blocking func
}
