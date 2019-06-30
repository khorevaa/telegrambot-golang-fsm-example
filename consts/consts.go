// constants
package consts

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

const (
	ItemChoose              string = "STATE:ITEM_CHOOSE"
	Name                    string = "STATE:NAME"
	Gender                  string = "STATE:GENDER"
	StatsCallbackDataPrefix string = "ukey"
)

var ItemsKeyboard = [][]tgbotapi.InlineKeyboardButton{
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

var GendersMarkup = tgbotapi.ReplyKeyboardMarkup{
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
