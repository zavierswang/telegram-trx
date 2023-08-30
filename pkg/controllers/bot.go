package controllers

import "github.com/mr-linch/go-tg"

var Menu = struct {
	Start    string
	Exchange string
	Rate     string
	OTC      string
	Advance  string
	Help     string
	Houbi    string
}{
	Start:    "🥳 开始",
	Exchange: "💵在线兑换",
	Rate:     "📈TRX汇率",
	OTC:      "🌎OTC汇率",
	Advance:  "💰TRX预支",
	Help:     "☎️咨询客服",
	Houbi:    "🔥火币交易",
}

type Bot struct {
	ReplayMarkup *tg.ReplyKeyboardMarkup
	Cmd          []tg.BotCommand
}

func NewBot() *Bot {
	layout := tg.NewReplyKeyboardMarkup(
		tg.NewButtonRow(
			tg.NewKeyboardButton(Menu.Exchange),
			tg.NewKeyboardButton(Menu.Advance),
		),
		tg.NewButtonRow(
			tg.NewKeyboardButton(Menu.OTC),
			tg.NewKeyboardButton(Menu.Rate),
		),
	)
	layout.ResizeKeyboard = true

	botCmd := []tg.BotCommand{
		{Command: "start", Description: Menu.Start},
		{Command: "rate", Description: Menu.Rate},
	}
	return &Bot{
		ReplayMarkup: layout,
		Cmd:          botCmd,
	}
}
