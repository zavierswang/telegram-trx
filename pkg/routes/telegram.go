package routes

import (
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"regexp"
	"telegram-trx/pkg/controllers"
	"telegram-trx/pkg/middleware"
)

func Telegram(router *tgb.Router) {
	router.Use(middleware.SessionManager)
	router.Use(tgb.MiddlewareFunc(middleware.Hook))

	router.Message(controllers.Start, tgb.Command("start"), tgb.ChatType(tg.ChatTypePrivate))
	router.Message(controllers.Rate, tgb.Any(tgb.Command("rate"), tgb.TextEqual(controllers.Menu.Rate)))
	router.Message(controllers.Exchange, tgb.TextEqual(controllers.Menu.Exchange), tgb.ChatType(tg.ChatTypePrivate))
	router.Message(controllers.Otc, tgb.TextEqual(controllers.Menu.OTC), tgb.ChatType(tg.ChatTypePrivate))
	router.Message(controllers.Advance, tgb.TextEqual(controllers.Menu.Advance), tgb.ChatType(tg.ChatTypePrivate))
	router.Message(controllers.AdvanceAddress, tgb.Any(tgb.Regexp(regexp.MustCompile(`^T\w+`)), middleware.IsSessionStep(middleware.SessionAdvance)), tgb.ChatType(tg.ChatTypePrivate))

	router.CallbackQuery(controllers.Retry, tgb.Regexp(regexp.MustCompile(`^retry\s+T\w+\s+([0-9]*[.])?([0-9]+)?`)))
}
