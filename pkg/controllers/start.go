package controllers

import (
	"context"
	"fmt"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"go.uber.org/zap/buffer"
	"html/template"
	"strconv"
	"strings"
	"telegram-trx/pkg/core/cst"
	"telegram-trx/pkg/core/global"
	"telegram-trx/pkg/core/logger"
	"telegram-trx/pkg/services"
)

type StartTmpl struct {
	Support string
	Px      string
	Px10    string
	Address string
}

func Start(ctx context.Context, update *tgb.MessageUpdate) error {
	userId := update.Message.From.ID.PeerID()
	username := update.Message.From.Username.PeerID()
	username = strings.ReplaceAll(username, "@", "")
	logger.Info("[%s %s] trigger action [start] controller", userId, username)
	bot := NewBot()
	err := update.Client.SetMyCommands(bot.Cmd).DoVoid(ctx)
	if err != nil {
		logger.Error("[%s %s] set command failed %v", userId, username, err)
		return err
	}
	market, err := services.GetOkxMarketRate()
	if err != nil {
		logger.Error("[%s %s] services.GetOkxMarketRate() failed %v", userId, username, err)
		return update.Answer("网络错误，请稍后~").DoVoid(ctx)
	}
	var _rate services.Rate
	for _, r := range market.Data {
		if r.Side == "buy" {
			_rate = r
			break
		}
	}
	__rate, _ := strconv.ParseFloat(_rate.Px, 64)
	_frate := 1 / __rate * (1.0 - global.App.Config.Telegram.Ratio)

	buf := new(buffer.Buffer)
	tmpl, err := template.ParseFiles(cst.StartTemplateFile)
	if err != nil {
		logger.Error("[%s %s] template parse file %s, failed %v", userId, username, cst.StartTemplateFile, err)
		return err
	}
	tpl := StartTmpl{
		Support: global.App.Config.App.Support,
		Px:      fmt.Sprintf("%.3f", _frate),
		Px10:    fmt.Sprintf("%.3f", _frate*10),
		Address: global.App.Config.Telegram.ReceiveAddress,
	}
	if err := tmpl.Execute(buf, tpl); err != nil {
		logger.Error("[%s %s] template execute file %s, failed %v", userId, username, cst.StartTemplateFile, err)
		return err
	}
	return update.Answer(buf.String()).ParseMode(tg.HTML).ReplyMarkup(bot.ReplayMarkup).DoVoid(ctx)
}
