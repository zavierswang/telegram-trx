package bootstrap

import (
	"context"
	"fmt"
	"github.com/mr-linch/go-tg"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap/buffer"
	"html/template"
	"math"
	"os"
	"strconv"
	"telegram-trx/pkg/core/cst"
	"telegram-trx/pkg/core/global"
	"telegram-trx/pkg/core/logger"
	"telegram-trx/pkg/models"
	"telegram-trx/pkg/services"
	"telegram-trx/pkg/services/grid"
	"telegram-trx/pkg/services/tron"
	"telegram-trx/pkg/utils"
	"time"
)

var (
	params  map[string]string
	headers map[string]string
	RatePx  string
	Px      float64
)

func StartCron() {
	global.App.Cron = cron.New(cron.WithSeconds())
	go func() {
		listenUSDT := &ListenUSDT{
			ctx: context.Background(),
		}
		_, err := global.App.Cron.AddJob("*/30 * * * * *", listenUSDT)
		if err != nil {
			logger.Error("[scheduler] add job listenUSDT failed %v", err)
			return
		}
		global.App.Cron.Start()
		defer global.App.Cron.Stop()
		select {}
	}()
}

type ListenUSDT struct {
	ctx        context.Context
	ticker     int
	trc20Queue []string
}

func (l *ListenUSDT) Run() {
	l.ticker++
	//logger.Info("[scheduler] run %d times, starting", l.ticker)
	now := time.Now()
	params = map[string]string{
		"limit":            "30",
		"start":            "0",
		"relatedAddress":   global.App.Config.Telegram.ReceiveAddress,
		"contract_address": cst.ContractAddress,
		"sort":             "-timestamp",
		"count":            "true",
		"filterTokenValue": "0",
		"start_timestamp":  strconv.FormatInt(now.Add(-120*time.Second).UnixMilli(), 10),
		"end_timestamp":    strconv.FormatInt(now.UnixMilli(), 10),
	}
	headers = map[string]string{
		"TRON-PRO-API-KEY": global.App.Config.Telegram.TronScanApiKey,
	}
	transfers, err := tron.TRC20Transfer(params, headers, true, true)
	if err != nil {
		logger.Error("[scheduler] tron.TRC20Transfer failed %v", err)
		return
	}
	// 首次启动并获取到的数据暂存到队列中
	if l.ticker == 1 {
		for _, transfer := range transfers {
			l.trc20Queue = append(l.trc20Queue, transfer.TransactionId)
		}
		//logger.Info("scheduler ListenUSDT %d times, trc20 %+v", l.ticker, l.trc20Queue)
		return
	}
	slice1 := l.trc20Queue
	var slice2 []string
	for _, transfer := range transfers {
		slice2 = append(slice2, transfer.TransactionId)
	}
	// 比对历史数据获取最新的交易号
	txIds, _ := utils.Comp(slice1, slice2)
	if len(txIds) != 0 {
		logger.Info("[scheduler] latest txIds coming, will be get market rate.")
		market, err := services.GetOkxMarketRate()
		if err != nil || market == nil {
			logger.Warn("[scheduler] services.GetOkxMarketRate failed %v", err)
			//interfaceRate, exist := global.Cache.Get("rate")
			//if !exist {
			//	logger.Warn("not found rate in cache.")
			//	return
			//}
			//RatePx = interfaceRate.(string)
			//logger.Info("[scheduler] rate in cache found, %s", RatePx)
		}
		RatePx = market.Data[0].Px
		logger.Info("[scheduler] market rate will be cached")
		//global.Cache.Set("rate", ratePx, cache.NoExpiration)
		Px, err = strconv.ParseFloat(RatePx, 64)
		if err != nil {
			logger.Error("[scheduler] market rate %s, failed %v", RatePx, err)
			return
		}
	}

	for _, txId := range txIds {
		for _, transfer := range transfers {
			if transfer.TransactionId == txId {
				l.exec(transfer, Px)
			}
		}
	}
	// 历史数据过多，删除部分数据
	if len(l.trc20Queue) >= 500 {
		logger.Info("[scheduler] clean remain queue ...")
		l.trc20Queue = l.trc20Queue[400:]
	}
	// 合并历史交易号，获新的数据
	l.trc20Queue = utils.Union(l.trc20Queue, slice2)
	return
}

func (l *ListenUSDT) exec(transfer tron.Transfer, Px float64) {
	var quant float64
	var balance float64
	txId := transfer.TransactionId
	quant, err := strconv.ParseFloat(transfer.Quant, 64)
	if err != nil {
		logger.Error("[scheduler] strconv.ParseFloat failed %v", err)
		return
	}
	//过滤小于0.01USDT的转帐
	if quant <= math.Pow10(4) {
		logger.Warn("[scheduler] %s is too small amount will be ignore.", txId)
		return
	}
	balance = utils.Trunc(quant/math.Pow10(6), 3)
	logger.Info("[scheduler] ListenUSDT %d times, found latest txId: %s, amount: %.3f USDT from %s", l.ticker, txId, balance, transfer.FromAddress)
	createTime := time.UnixMilli(transfer.BlockTs).Format(cst.DateTimeFormatter)
	usdt := USDT{
		TransactionId: transfer.TransactionId,
		FromAddress:   transfer.FromAddress,
		ToAddress:     transfer.ToAddress,
		CreateTime:    createTime,
		Balance:       balance,
	}
	logger.Info("[scheduler] transfer data => %+v", usdt)
	err = l.transfer(usdt, Px)
	if err != nil {
		logger.Error("[scheduler] transfer trx %+v failed %v", usdt, err)
		return
	}
	return
}

func (l *ListenUSDT) transfer(usdt USDT, Px float64) error {
	groupId := tg.Username(global.App.Config.App.Support)
	ray := 1.0 / Px
	balance := ray * usdt.Balance
	amount := balance * (1.0 - global.App.Config.Telegram.Ratio)
	amount = utils.Trunc(amount, 3)
	order := models.Order{
		TxID:        usdt.TransactionId,
		FromAddress: usdt.FromAddress,
		ToAddress:   usdt.ToAddress,
		Balance:     usdt.Balance,
		Amount:      usdt.Amount,
		Status:      cst.OrderStatusRunning,
		Finished:    false,
	}
	err := global.App.DB.Save(&order).Error
	if err != nil {
		logger.Error("[scheduler] new order record failed %v", err)
		usdt.Status = "失败"
		usdt.Description = "服务器异常"
		l.failure(usdt)
		return nil
	}
	logger.Info("[scheduler] transfer amount %.3f TRX", amount)
	// 查看是否有预支记录
	addr := models.Address{
		Address: usdt.FromAddress,
		Balance: 0.0,
		Advance: 0.0,
		Count:   0,
	}
	global.App.DB.FirstOrCreate(&addr, "address = ?", usdt.FromAddress)
	// 扣除预支金额
	amount -= addr.Advance
	// 当前兑换不足以抵扣预支金额
	if amount <= 0 {
		addr.Advance = math.Abs(amount)
		logger.Warn("[scheduler] exchange amount is not sufficient to offset the advance")
		logger.Warn("[scheduler] still remaining of advance amount %.3f", addr.Advance)
		global.App.DB.Save(&addr)
		usdt.Status = "成功"
		usdt.Description = "不足以抵扣预支金额"
		l.failure(usdt)
		return nil
	}
	count := int(math.Floor(usdt.Balance / global.App.Config.Telegram.ThresholdValue))
	logger.Info("[scheduler] %s increment %d count of advance", usdt.FromAddress, count)
	addr.Count += count
	addr.Balance += usdt.Balance
	global.App.DB.Save(&addr)
	logger.Info("[scheduler] %s deduct the advance and pay %.3f TRX", usdt.FromAddress, amount)
	logger.Info("[scheduler] will be transfer %s => %f TRX", usdt.FromAddress, amount)
	// 还需支付金额
	usdt.Amount = amount
	//检查余额是否足够
	asset, err := tron.TokenAssetOverview(global.App.Config.Telegram.SendAddress)
	if err != nil {
		logger.Error("[scheduler] tron.TokenAssetOverview failed %v", err)
		usdt.Status = "失败"
		usdt.Description = "网络错误"
		return err
	}
	var remainAmount float64 // 钱包余额
	if asset != nil {
		b, _ := strconv.ParseFloat(asset.Balance, 64)
		remainAmount = utils.Trunc(b/math.Pow10(6), 2)
	} else {
		remainAmount = 0.0
	}
	logger.Info("[scheduler] account out wallet balance %.3f, %3f TRX will be transferred to %s", remainAmount, amount, usdt.FromAddress)
	if remainAmount < amount {
		tmpl, _ := template.ParseFiles(cst.RemainTemplateFile)
		buf := new(buffer.Buffer)
		tpl := AssetSufficient{
			usdt,
			remainAmount,
			amount,
		}
		_ = tmpl.Execute(buf, tpl)
		global.App.DB.Model(&models.Order{}).Where("tx_id = ?", usdt.TransactionId).Updates(map[string]interface{}{"status": cst.OrderStatusNotSufficientFunds, "finished": true})
		logger.Warn("[scheduler] %s not sufficient funds, remain: %.3f, will be paid: %.3f", global.App.Config.Telegram.SendAddress, remainAmount, amount)
		layout := tg.NewButtonLayout[tg.InlineKeyboardButton](1).Row()
		layout.Insert(tg.NewInlineKeyboardButtonCallback("确认钱包余额充足，手动补充订单", fmt.Sprintf("retry %s %f", usdt.TransactionId, amount)))
		inlineKeyboard := tg.NewInlineKeyboardMarkup(layout.Keyboard()...)
		_ = global.App.Client.SendMessage(groupId, buf.String()).ParseMode(tg.HTML).ReplyMarkup(inlineKeyboard).DoVoid(context.Background())
		usdt.Status = "失败"
		usdt.Description = "余额不足，请联系管理补发"
		l.failure(usdt)
		return nil
	}

	logger.Info("[scheduler] start transfer USDT: %.3f => TRX: %.3f", usdt.Balance, usdt.Amount)
	txID, err := grid.TransferTRX(global.App.Config.Telegram.SendAddress, usdt.FromAddress, usdt.Amount)
	if err != nil {
		logger.Error("[scheduler] grid transfer failed %v", err)
		order.Finished = true
		order.Status = cst.OrderStatusApiFailure
		usdt.Status = "失败"
		usdt.Description = "转帐失败"
		global.App.DB.Save(&order)
		l.failure(usdt)
		return err
	}
	logger.Info("[scheduler] grid transfer successfully, txId: %s", txID)

	order.Status = cst.OrderStatusSuccess
	order.Finished = true
	usdt.Status = "成功"
	global.App.DB.Save(&order)
	logger.Info("[scheduler] transfer USDT: %.2f ==> %s TRX: %.2f successfully", usdt.Balance, usdt.FromAddress, usdt.Amount)
	l.success(usdt)
	return nil
}

func (l *ListenUSDT) success(usdt USDT) {
	groups := global.App.Config.App.Groups
	pf, _ := os.ReadFile(cst.SuccessTemplateFile)
	tmpl, err := template.New("success").Parse(string(pf))
	if err != nil {
		logger.Error("[scheduler] template parse file %s, failed %v", cst.SuccessTemplateFile, err)
		return
	}
	buf := new(buffer.Buffer)
	err = tmpl.Execute(buf, usdt)
	if err != nil {
		logger.Error("[scheduler] template execute file %s, failed %v", cst.SuccessTemplateFile, err)
		return
	}
	for _, group := range groups {
		groupId := tg.Username(group)
		err = global.App.Client.SendMessage(groupId, buf.String()).ParseMode(tg.HTML).DoVoid(l.ctx)
		if err != nil {
			logger.Error("[scheduler] send message to group failed %v", err)
		}
	}
}

func (l *ListenUSDT) failure(usdt USDT) {
	groups := global.App.Config.App.Groups
	pf, _ := os.ReadFile(cst.FailureTemplateFile)
	tmpl, _ := template.New("failure").Funcs(template.FuncMap{}).Parse(string(pf))
	buf := new(buffer.Buffer)
	_ = tmpl.Execute(buf, usdt)
	for _, group := range groups {
		groupId := tg.Username(group)
		_ = global.App.Client.SendMessage(groupId, buf.String()).ParseMode(tg.HTML).DoVoid(l.ctx)
	}
}

type USDT struct {
	TransactionId string
	FromAddress   string
	ToAddress     string
	CreateTime    string
	Balance       float64
	Status        string
	Description   string
	Amount        float64
}

type AssetSufficient struct {
	USDT
	Remaining float64
	Amount    float64
}
