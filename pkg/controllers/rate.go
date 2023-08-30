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
	"time"
)

type RealTimeRate struct {
	Px      string `json:"px"`
	Px5     string `json:"px5"`
	Px10    string `json:"px10"`
	Px20    string `json:"px20"`
	Px50    string `json:"px50"`
	Px100   string `json:"px100"`
	Px500   string `json:"px500"`
	Px1000  string `json:"px1000"`
	Time    string `json:"time"`
	Address string `json:"address"`
}

func Rate(ctx context.Context, update *tgb.MessageUpdate) error {
	userId := update.Message.From.ID.PeerID()
	username := update.Message.From.Username.PeerID()
	username = strings.ReplaceAll(username, "@", "")
	logger.Info("[%s %s] trigger action [rate] controller", userId, username)
	market, err := services.GetOkxMarketRate()
	if err != nil {
		logger.Error("[%s %s] services.GetOkxMarketRate failed %v", userId, username, err)
		return update.Answer("数据获取失败，请重试~").DoVoid(ctx)
	}
	var _rate services.Rate
	for _, r := range market.Data {
		if r.Side == "buy" {
			_rate = r
			break
		}
	}
	ts, _ := strconv.ParseInt(_rate.Ts, 10, 64)
	dateTime := time.UnixMilli(ts).Format(cst.DateTimeFormatter)
	__rate, _ := strconv.ParseFloat(_rate.Px, 64)

	_frate := 1 / __rate * (1.0 - global.App.Config.Telegram.Ratio)
	data := RealTimeRate{
		Px:      fmt.Sprintf("%.3f", _frate),
		Px5:     fmt.Sprintf("%.3f", _frate*5),
		Px10:    fmt.Sprintf("%.3f", _frate*10),
		Px20:    fmt.Sprintf("%.3f", _frate*20),
		Px50:    fmt.Sprintf("%.3f", _frate*50),
		Px100:   fmt.Sprintf("%.3f", _frate*100),
		Px500:   fmt.Sprintf("%.3f", _frate*500),
		Px1000:  fmt.Sprintf("%.3f", _frate*1000),
		Time:    dateTime,
		Address: global.App.Config.Telegram.ReceiveAddress,
	}
	b := new(buffer.Buffer)
	tmpl, err := template.ParseFiles(cst.RateTemplateFile)
	if err != nil {
		logger.Error("[%s %s] template parse file %s, failed %v", userId, username, cst.RateTemplateFile, err)
		return err
	}
	err = tmpl.Execute(b, data)
	if err != nil {
		logger.Error("[%s %s] template execute failed %v", userId, username, err)
		return err
	}
	return update.Answer(b.String()).ParseMode(tg.HTML).DisableWebPagePreview(true).DoVoid(ctx)
}
