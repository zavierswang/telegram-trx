package controllers

import (
	"context"
	"github.com/mr-linch/go-tg/tgb"
	"regexp"
	"strconv"
	"strings"
	"telegram-trx/pkg/core/cst"
	"telegram-trx/pkg/core/global"
	"telegram-trx/pkg/core/logger"
	"telegram-trx/pkg/models"
	"telegram-trx/pkg/services/grid"
	"telegram-trx/pkg/utils"
)

func Retry(ctx context.Context, callback *tgb.CallbackQueryUpdate) error {
	chatId := callback.Message.Chat.ID
	userId := callback.Message.From.ID.PeerID()
	username := callback.Message.From.Username.PeerID()
	username = strings.ReplaceAll(username, "@", "")
	logger.Info("[%s %s] trigger action [retry] callback", userId, username)
	//messageId := callback.Message.ID
	compile, err := regexp.Compile(`^retry\s+(?P<txID>T\w+)\s+(?P<amount>[0-9]*[.]?[0-9]+)`)
	if err != nil {
		logger.Error("[%s %s] compile failed %v", userId, username, err)
		return err
	}
	groups := utils.FindGroups(compile, callback.Data)
	var users []models.User
	global.App.DB.Find(&users, "is_admin = ? AND user_id = ?", true, userId)
	if len(users) == 0 {
		logger.Warn("[%s %s] %s is not administrator", userId, username, userId)
		return nil
	}
	txID := groups["txID"]
	amount, _ := strconv.ParseFloat(groups["amount"], 64)
	var orders []models.Order
	global.App.DB.Find(&orders, "tx_id = ? AND amount = ? AND finished = ?", txID, amount, false)
	if len(orders) == 0 {
		logger.Warn("[%s %s] not found order", userId, username)
		return nil
	}
	_ = callback.Client.SendMessage(chatId, "正在处理~").
		DisableNotification(false).
		DoVoid(ctx)

	toAddress := orders[0].FromAddress
	txId, err := grid.TransferTRX(global.App.Config.Telegram.SendAddress, toAddress, amount)
	if err != nil {
		logger.Error("[%s %s] retry transfer TRX failed %v", userId, username, err)
		return callback.Client.SendMessage(chatId, "手动补单失败，请重新确认钱包余额或联系开发人员").DisableNotification(false).DoVoid(ctx)
	}
	logger.Info("[%s %s] grid transfer trx successfully, txId: %s", txId)
	order := orders[0]
	order.Status = cst.OrderStatusSuccess
	order.Finished = true
	global.App.DB.Save(&order)
	return callback.Client.SendMessage(chatId, "处理成功~").
		DisableNotification(false).
		DoVoid(ctx)

}
