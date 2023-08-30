package controllers

import (
	"context"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"go.uber.org/zap/buffer"
	"html/template"
	"os"
	"strings"
	"telegram-trx/pkg/core/cst"
	"telegram-trx/pkg/core/logger"
	"telegram-trx/pkg/services"
	"telegram-trx/pkg/utils"
	"time"
)

func Otc(ctx context.Context, update *tgb.MessageUpdate) error {
	userId := update.Message.From.ID.PeerID()
	username := update.Message.From.Username.PeerID()
	username = strings.ReplaceAll(username, "@", "")
	logger.Info("[%s %s] trigger action [rate] controller", userId, username)
	trading, err := services.GetOkxTradingOrders()
	if err != nil {
		logger.Error("[%s %s] services.GetOkxTradingOrders failed %v", userId, username, err)
		return update.Answer("数据获取失败，请重试~").DoVoid(ctx)
	}
	sells := trading.Data.Sell
	if len(trading.Data.Sell) >= 10 {
		sells = trading.Data.Sell[:10]
	}
	pf, err := os.ReadFile(cst.OtcTemplateFile)
	if err != nil {
		logger.Error("[%s %s] template read file %s, failed %v", userId, username, cst.OtcTemplateFile, err)
		return err
	}
	var res = struct {
		Sells    []services.RateSell
		DateTime time.Time
	}{
		Sells:    sells,
		DateTime: time.Now(),
	}
	tmpl, err := template.New("otc").Funcs(template.FuncMap{
		"add":    utils.AddIndex,
		"format": utils.FormatTime,
	}).Parse(string(pf))
	buf := new(buffer.Buffer)
	err = tmpl.Execute(buf, res)
	if err != nil {
		logger.Error("[%s %s] template execute file %s, failed %v", userId, username, cst.OtcTemplateFile, err)
		return err
	}
	return update.Answer(buf.String()).ParseMode(tg.HTML).DoVoid(ctx)
}
