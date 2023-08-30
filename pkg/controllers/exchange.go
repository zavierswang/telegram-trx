package controllers

import (
	"context"
	"fmt"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"go.uber.org/zap/buffer"
	"html/template"
	"math"
	"os"
	"strconv"
	"strings"
	"telegram-trx/pkg/core/cst"
	"telegram-trx/pkg/core/global"
	"telegram-trx/pkg/core/logger"
	"telegram-trx/pkg/middleware"
	"telegram-trx/pkg/models"
	"telegram-trx/pkg/services/tron"
)

type ExchangeTmpl struct {
	Balance string
	Address string
}

func Exchange(ctx context.Context, update *tgb.MessageUpdate) error {
	userId := update.Message.From.ID.PeerID()
	username := update.Message.From.Username.PeerID()
	username = strings.ReplaceAll(username, "@", "")
	logger.Info("[%s %s] trigger action [exchange] controller", userId, username)
	sess := middleware.SessionManager.Get(ctx)
	middleware.SessionManager.Reset(sess)
	uid := update.From.ID.PeerID()
	var user models.User
	var balance float64
	err := global.App.DB.First(&user, "user_id = ?", uid).Error
	if err != nil {
		logger.Error("[%s %s] not found user", userId, username)
		return update.Answer("非法用户~").DoVoid(ctx)
	}
	if user.IsAdmin {
		asset, err := tron.TokenAssetOverview(global.App.Config.Telegram.SendAddress)
		if err != nil {
			logger.Error("[%s %s] tron.TokenAssetOverview failed %v", userId, username, err)
			return err
		}
		logger.Info("[%s %s] token asset overview: %+v", userId, username, asset)
		if asset == nil {
			logger.Warn("[%s %s] token asset overview is nil", userId, username)
			balance = 0.0
		} else {
			logger.Info("[%s %s] token asset overview: %+v", userId, username, asset)
			f, err := strconv.ParseFloat(asset.Balance, 64)
			if err != nil {
				logger.Error("[%s %s] asset balance parse float failed %v", userId, username, err)
				balance = 0.0
			}
			balance = f / math.Pow10(6)
		}
	} else {
		balance = 0.0
	}
	tpl := ExchangeTmpl{
		Balance: fmt.Sprintf("%.3f", balance),
		Address: global.App.Config.Telegram.ReceiveAddress,
	}
	buf := new(buffer.Buffer)
	pf, _ := os.ReadFile(cst.ExchangeTemplateFile)
	tmpl, err := template.New("exchange").Parse(string(pf))
	if err != nil {
		logger.Error("[%s %s] %s failed %v", userId, username, cst.ExchangeTemplateFile, err)
		return err
	}
	err = tmpl.Execute(buf, tpl)
	if err != nil {
		logger.Error("template execute failed %v", err)
		return err
	}
	layout := tg.NewButtonLayout[tg.InlineKeyboardButton](2).Row()
	layout.Insert(tg.NewInlineKeyboardButtonURL("联系客服👩‍💻", fmt.Sprintf("https://t.me/%s", global.App.Config.App.Support)))
	inlineKeyboard := tg.NewInlineKeyboardMarkup(layout.Keyboard()...)

	var inputFile tg.InputFile
	receiverAddressIconPath := global.App.Config.Telegram.ReceiveAddressIcon
	if receiverAddressIconPath != "" {
		inputFile, err = tg.NewInputFileLocal(receiverAddressIconPath)
		if err != nil {
			logger.Error("[%s %s] load receiver icon %s failed %v", userId, username, receiverAddressIconPath, err)
			return err
		}
		fileArg := tg.NewFileArgUpload(inputFile)
		return update.AnswerPhoto(fileArg).Caption(buf.String()).ParseMode(tg.HTML).ReplyMarkup(inlineKeyboard).DoVoid(ctx)
	}
	return update.Answer(buf.String()).ParseMode(tg.HTML).ReplyMarkup(inlineKeyboard).DoVoid(ctx)
}
