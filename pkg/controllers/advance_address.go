package controllers

import (
	"context"
	"fmt"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"go.uber.org/zap/buffer"
	"html/template"
	"strings"
	"telegram-trx/pkg/core/cst"
	"telegram-trx/pkg/core/global"
	"telegram-trx/pkg/core/logger"
	"telegram-trx/pkg/middleware"
	"telegram-trx/pkg/models"
	"telegram-trx/pkg/services/grid"
	"time"
)

func AdvanceAddress(ctx context.Context, update *tgb.MessageUpdate) error {
	userId := update.Message.From.ID.PeerID()
	username := update.Message.From.Username.PeerID()
	username = strings.ReplaceAll(username, "@", "")
	logger.Info("[%s %s] trigger action [address] controller", userId, username)
	valid, err := grid.ValidateAddress(update.Text)
	if err != nil || !valid.Result {
		logger.Error("[%s %s] validate wallet address failed %v", userId, username, err)
		return update.Answer("非法地址，请输入TRC20地址").DoVoid(ctx)
	}
	sess := middleware.SessionManager.Get(ctx)
	sess.Address = update.Text
	walletAddress := update.Text
	//重置session
	middleware.SessionManager.Reset(sess)
	var tmpl *template.Template
	layout := tg.NewButtonLayout[tg.InlineKeyboardButton](2).Row()
	layout.Insert(
		tg.NewInlineKeyboardButtonURL("联系管理👩‍💻", fmt.Sprintf("https://t.me/%s", global.App.Config.App.Support)),
	)
	inlineKeyboard := tg.NewInlineKeyboardMarkup(layout.Keyboard()...)

	var addrs []models.Address
	global.App.DB.Find(&addrs, "address = ?", walletAddress)
	if len(addrs) == 0 {
		buf := new(buffer.Buffer)
		tmpl, _ = template.ParseFiles(cst.AdvanceFailureTemplateFile)
		_ = tmpl.Execute(buf, nil)
		return update.Update.Reply(ctx, update.Answer(buf.String()).
			ParseMode(tg.HTML).
			DisableWebPagePreview(false).
			ReplyMarkup(inlineKeyboard),
		)
	}
	addr := addrs[0]
	logger.Info("[%s %s] compute address %+v advance permission", userId, username, addr)
	if addr.Count <= 0 || addr.Advance > 0 {
		logger.Error("[%s %s] forbiden advance", userId, username)
		buf := new(buffer.Buffer)
		tmpl, _ = template.ParseFiles(cst.AdvanceFailureTemplateFile)
		_ = tmpl.Execute(buf, nil)
		return update.Update.Reply(ctx, update.Answer(buf.String()).
			ParseMode(tg.HTML).
			DisableWebPagePreview(false).
			ReplyMarkup(inlineKeyboard),
		)
	}
	logger.Warn("[%s %s] accept advance to address %s", userId, username, walletAddress)
	buf := new(buffer.Buffer)
	addr.Advance += global.App.Config.Telegram.AdvanceAmount
	addr.Count -= 1
	//groupId := tg.Username(global.App.Config.App.Support)
	advanceAmount := global.App.Config.Telegram.AdvanceAmount
	txId, err := grid.TransferTRX(global.App.Config.Telegram.SendAddress, walletAddress, advanceAmount)
	if err != nil {
		logger.Error("[%s %s] advance to %s failed %v", userId, username, walletAddress, err)
		return err
	}
	logger.Info("[%s %s] advance to %s successfully，amount %.3f TRX, txID: %s", userId, username, walletAddress, advanceAmount, txId)

	global.App.DB.Model(&models.Address{}).
		Where("address = ?", walletAddress).
		Updates(map[string]interface{}{"advance": addr.Advance, "count": addr.Count})
	// 群通知预支记录
	t, _ := template.ParseFiles(cst.AdvanceGroupTemplateFile)
	b := new(buffer.Buffer)
	d := struct {
		Address  string
		Datetime string
	}{
		Address:  walletAddress,
		Datetime: time.Now().Format(cst.DateTimeFormatter),
	}
	_ = t.Execute(b, d)
	groups := global.App.Config.App.Groups
	for _, group := range groups {
		groupId := tg.Username(group)
		_ = global.App.Client.SendMessage(groupId, b.String()).ParseMode(tg.HTML).DoVoid(ctx)
	}
	tmpl, _ = template.ParseFiles(cst.AdvanceSuccessTemplateFile)
	_ = tmpl.Execute(buf, nil)
	return update.Update.Reply(ctx, update.Answer(buf.String()).
		ParseMode(tg.HTML).
		DisableWebPagePreview(true),
	)
}
