package callbacks

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"fsmExample/consts"
	"fsmExample/database"
)

// handle Emoji clicks
func ClickOnItem(api tgbotapi.BotAPI, update tgbotapi.Update) {
	text := fmt.Sprintf("This is %s", update.CallbackQuery.Data)
	action := tgbotapi.NewCallback(update.CallbackQuery.ID, text)
	_, _ = api.AnswerCallbackQuery(action)

	edit := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    update.CallbackQuery.Message.Chat.ID,
			MessageID: update.CallbackQuery.Message.MessageID,
		},
		Text:      HtmlFmt("Now please enter your pretty name:", "b"),
		ParseMode: "html",
	}

	_, _ = api.Send(edit)

	database.UpdateState(update, consts.Name)
	database.UpdateData(update, map[string]interface{}{"fav_item": update.CallbackQuery.Data})
}

// process after-/stats click
func ProcessAbout(api tgbotapi.BotAPI, update tgbotapi.Update) {
	info := database.Get(strings.Replace(update.CallbackQuery.Data, "ukey:", "", 1), false)
	savedDataString := fmt.Sprintf(
		"Name: %s\nGender: %s\nFavorite item: %s\nAge ~ %s",
		fmt.Sprintf("%v", info["name"]),
		fmt.Sprintf("%v", info["gender"]),
		fmt.Sprintf("%v", info["fav_item"]),
		fmt.Sprintf("%v", info["age"]),
	)

	action := tgbotapi.NewCallbackWithAlert(update.CallbackQuery.ID, savedDataString)
	_, _ = api.AnswerCallbackQuery(action)
}
