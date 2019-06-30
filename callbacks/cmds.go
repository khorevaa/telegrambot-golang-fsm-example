// handler for slash-commands

// btw, that's a bad idea to store text messages in code

package callbacks

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"fsmExample/consts"
	"fsmExample/database"
)

// process /start
func CmdStart(api tgbotapi.BotAPI, update tgbotapi.Update) {
	buttons := [][]tgbotapi.InlineKeyboardButton{
		{
			tgbotapi.NewInlineKeyboardButtonData("🔧", "toggle"),
			tgbotapi.NewInlineKeyboardButtonData("🧱", "bricks"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("⚒", "forge"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("📸", "camera"),
			tgbotapi.NewInlineKeyboardButtonData("📞", "telescope"),
			tgbotapi.NewInlineKeyboardButtonData("☎️", "phone"),
		},
	}

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		HtmlFmt("We'll try to predict you age. But first choose thing you want:", "b"))
	msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
	msg.ParseMode = "html"

	_, _ = api.Send(msg)

	database.UpdateState(update, consts.ItemChoose)
}

func CmdStats(api tgbotapi.BotAPI, update tgbotapi.Update) {
	all := database.GetAllKeys(10)

	msg := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		HtmlFmt("Here are all users, who participated:", "b"))
	msg.ParseMode = "html"

	if len(all) == 0 {
		msg.Text = HtmlFmt("Oops! No one has participated...", "b")
	}

	var buttons [][]tgbotapi.InlineKeyboardButton

	for key := range all {
		button := []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s) %s", strconv.FormatInt(int64(key+1), 10), string(strings.Split(all[key], ":")[2])),
			fmt.Sprintf("ukey:%s", string(all[key])))}

		buttons = append(buttons, button)
	}

	msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: buttons,
	}
	_, _ = api.Send(msg)
}
