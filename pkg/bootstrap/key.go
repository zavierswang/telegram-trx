package bootstrap

import (
	"os"
	"telegram-trx/pkg/core/global"
	"telegram-trx/pkg/core/logger"
	"telegram-trx/pkg/services/grid"
)

func StoreKey() {
	flag := false
	accounts := grid.DescribeLocalAccounts()
	for key, ks := range accounts {
		for _, k := range ks {
			logger.Info("ks: %+v", k)
			if k.Address.String() == global.App.Config.Telegram.SendAddress {
				flag = true
				logger.Info("%s => %s has been stored", key, k.Address.String())
				return
			}
		}
	}
	if !flag {
		err := grid.ImportPrivateKey(global.App.Config.Telegram.PrivateKey, global.App.Config.Telegram.AliasKey)
		if err != nil {
			flag = false
			logger.Error("import private key failed %v", err)
		}
		return
	}
	if !flag {
		os.Exit(-1)
	}
}
