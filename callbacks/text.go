package callbacks

import (
	"fmt"
	"html"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"fsmExample/consts"
	"fsmExample/database"
)

// process
func ProcessName(api tgbotapi.BotAPI, update tgbotapi.Update) {
	markup := tgbotapi.ReplyKeyboardMarkup{
		Selective:       true,
		OneTimeKeyboard: true,
		ResizeKeyboard:  true,
		Keyboard: [][]tgbotapi.KeyboardButton{
			{
				{"🕺 " + "Male", false, false},
				{"🙋🏻‍♀️ " + "Female", false, false},
			},
		},
	}

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		fmt.Sprintf("<b> %s</b>, please enter your gender?", html.EscapeString(update.Message.Text)))
	msg.ReplyMarkup = markup
	msg.ParseMode = "html"
	_, _ = api.Send(msg)
	database.UpdateState(update, consts.Gender)
	database.UpdateData(update, map[string]interface{}{"name": update.Message.Text})
}

func ProcessGender(api tgbotapi.BotAPI, update tgbotapi.Update) {
	text := strings.ToLower(update.Message.Text)
	if !(strings.HasSuffix(text, "male") || strings.HasSuffix(text, "female")) {
		return
	}
	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		HtmlFmt("Predicting your age...", "code"))

	msg.ParseMode = "html"
	msg.ReplyMarkup = tgbotapi.ReplyKeyboardRemove{}

	message, err := api.Send(msg)

	if err != nil {
		log.Panic(err)
	}

	savedData := database.GetData(update)
	savedDataString := fmt.Sprintf(
		"Your name: %s\nYour gender: %s\nFavorite item: %s\n",
		HtmlFmt(fmt.Sprintf("%v", savedData["name"]), "b"),
		HtmlFmt(text, "code"),
		HtmlFmt(fmt.Sprintf("%v", savedData["fav_item"]), "i"))

	age := randInt(5, 40)
	edit := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    message.Chat.ID,
			MessageID: message.MessageID,
		},
		Text: fmt.Sprintf(
			"%sYour ~age is: %s",
			savedDataString,
			HtmlFmt(fmt.Sprintf("%d", age), "code")),
		ParseMode: "html",
	}

	time.Sleep(time.Second * 9)

	_, _ = api.Send(edit)
	database.UpdateData(update, map[string]interface{}{
		"gender": strings.Split(text, " ")[1],
		"age":    age,
	})
	database.UpdateState(update, database.InitialState)
}
