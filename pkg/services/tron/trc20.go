package tron

import (
	"fmt"
	"strconv"
	"telegram-trx/pkg/common/httpclient"
	"telegram-trx/pkg/core/cst"
	"telegram-trx/pkg/core/logger"
)

type TransfersResp struct {
	Total          int        `json:"total"`
	TokenTransfers []Transfer `json:"token_transfers"`
}

type Transfer struct {
	TransactionId   string    `json:"transaction_id"`
	BlockTs         int64     `json:"block_ts"`
	FromAddress     string    `json:"from_address"`
	ToAddress       string    `json:"to_address"`
	ContractAddress string    `json:"contract_address"`
	Quant           string    `json:"quant"`
	EventType       string    `json:"event_type"`
	ContractRet     string    `json:"contractRet"`
	ContractType    string    `json:"contract_type"`
	TokenInfo       TokenInfo `json:"tokenInfo"`
}

type NotificationTransfers struct {
	ID        int        `json:"id"`
	UserId    string     `json:"user_id"`
	Address   string     `json:"address"`
	Transfers []Transfer `json:"transfers"`
}

func TRC20Transfer(params map[string]string, headers map[string]string, doOnce bool, isCheckIn bool) ([]Transfer, error) {
	var resp TransfersResp
	url := fmt.Sprintf("%s%s", cst.TronScanApi, cst.TronTransfers)
	err := httpclient.GetJson(url, params, headers, &resp)
	if err != nil {
		logger.Error("[tron] %s request failed %v", url, err)
		return nil, err
	}
	//logger.Info("URL: %s", url)
	//logger.Info("Params: %+v", params)
	var transfers []Transfer
	for _, transfer := range resp.TokenTransfers {
		if transfer.ContractAddress == cst.ContractAddress && transfer.ContractRet == "SUCCESS" {
			//logger.Info("transfer: %+v", transfer.TransactionId)
			if isCheckIn {
				if transfer.ToAddress == params["relatedAddress"] {
					transfers = append(transfers, transfer)
				}
			} else {
				transfers = append(transfers, transfer)
			}
		}
	}
	if doOnce {
		return transfers, err
	}
	limit, _ := strconv.Atoi(params["limit"])
	page := 1

PAGE:
	if len(resp.TokenTransfers) >= limit {
		var body TransfersResp
		params["start"] = strconv.Itoa(page * limit)
		err = httpclient.GetJson(url, params, headers, &body)
		if err != nil {
			logger.Error("%s failed %v", url, err)
			return transfers, nil
		}
		page += 1
		if len(body.TokenTransfers) >= 30 {
			for _, transfer := range body.TokenTransfers {
				if transfer.ContractRet == "SUCCESS" && transfer.ContractAddress == cst.ContractAddress {
					if isCheckIn {
						if transfer.ToAddress == params["relatedAddress"] {
							transfers = append(transfers, transfer)
						}
					} else {
						transfers = append(transfers, transfer)
					}
				}
			}
			goto PAGE
		}
	}
	return transfers, nil
}
