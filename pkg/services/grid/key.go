package grid

import (
	"github.com/fbsobreira/gotron-sdk/pkg/account"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
	"github.com/fbsobreira/gotron-sdk/pkg/keystore"
	"github.com/fbsobreira/gotron-sdk/pkg/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"math"
	"telegram-trx/pkg/core/global"
	"telegram-trx/pkg/core/logger"
)

var (
	conn          *client.GrpcClient
	signerAddress tronAddress
	toAddress     tronAddress
)

const (
	defaultNode = "grpc.trongrid.io:50051"
	withTLS     = false
)

func ImportPrivateKey(privateKey string, accountName string) error {
	name, err := account.ImportFromPrivateKey(privateKey, accountName, "")
	if err != nil {
		return err
	}
	logger.Info("imported keystore given account alias of %s", name)
	add, _ := store.AddressFromAccountName(name)
	logger.Info("tron address: %s", add)
	return nil
}

func DescribeLocalAccounts() map[string][]keystore.Account {
	accounts := make(map[string][]keystore.Account, 1)
	for _, name := range store.LocalAccounts() {
		ks := store.FromAccountName(name)
		allAccounts := ks.Accounts()
		//logger.Info("%s => %+v", name, allAccounts)
		accounts[name] = allAccounts
	}
	return accounts
}

func TransferTRX(from, to string, amount float64) (string, error) {
	signerAddress = tronAddress{address: from}
	toAddress = tronAddress{address: to}
	value := int64(amount * math.Pow10(6))
	conn = client.NewGrpcClient(defaultNode)
	opts := make([]grpc.DialOption, 0)
	if withTLS {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(nil)))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	conn.SetAPIKey(global.App.Config.Telegram.GridApiKey)
	if err := conn.Start(opts...); err != nil {
		return "", err
	}
	tx, err := conn.Transfer(signerAddress.String(), toAddress.String(), value)
	if err != nil {
		logger.Error("transfer %s=> %s %d TRX failed %v", signerAddress.address, to, value, err)
		return "", err
	}
	logger.Info("txID: %s", common.Bytes2Hex(tx.GetTxid()))

	var ctrl *Controller
	ks, acct, err := store.UnlockedKeystore(signerAddress.String(), "")
	if err != nil {
		logger.Error("unlocked keystore failed %v", err)
		return "", err
	}
	ctrl = NewController(conn, ks, acct, tx.Transaction, func(controller *Controller) {})
	if err = ctrl.ExecuteTransaction(); err != nil {
		logger.Error("execute transaction failed %v", err)
		return "", err
	}
	return common.Bytes2Hex(tx.GetTxid()), nil
}

func findAddress(value string) (tronAddress, error) {
	address := tronAddress{}
	if err := address.Set(value); err != nil {
		acc, err := store.AddressFromAccountName(value)
		if err != nil {
			return address, err
		}
		return tronAddress{address: acc}, nil
	}
	return address, nil
}
