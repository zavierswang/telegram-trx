package tron

import (
	"fmt"
	"strconv"
	"telegram-trx/pkg/common/httpclient"
	"telegram-trx/pkg/core/cst"
	"telegram-trx/pkg/core/global"
	"telegram-trx/pkg/core/logger"
)

type Trc10Resp struct {
	Total      int         `json:"total"`
	RangeTotal int         `json:"rangeTotal"`
	Data       []Trc10Data `json:"data"`
}

type Trc10Data struct {
	Block           int64        `json:"block"`
	Hash            string       `json:"hash"`
	Timestamp       int64        `json:"timestamp"`
	OwnerAddress    string       `json:"ownerAddress"`
	ToAddress       string       `json:"toAddress"`
	ContractType    int          `json:"contractType"`
	Amount          string       `json:"amount"`
	ContractRet     string       `json:"contractRet"`
	ContractData    ContractData `json:"contractData"`
	TokenInfo       TokenInfo    `json:"tokenInfo"`
	TokenType       string       `json:"tokenType"`
	RiskTransaction bool         `json:"riskTransaction"`
}

type TokenInfo struct {
	TokenId   string `json:"tokenId"`
	TokenAbbr string `json:"tokenAbbr"`
	TokenName string `json:"tokenName"`
	TokenType string `json:"tokenType"`
}

type ContractData struct {
	Amount       int64  `json:"amount"`
	OwnerAddress string `json:"owner_address"`
	ToAddress    string `json:"to_address"`
}

func TRC10Transfer(params map[string]string, doOnce bool, isCheckIn bool) ([]Trc10Data, error) {
	headers := map[string]string{
		"TRON-PRO-API-KEY": global.App.Config.Telegram.TronScanApiKey,
	}
	var resp Trc10Resp
	url := fmt.Sprintf("%s%s", cst.TronScanApi, cst.TronTrc10Transfers)
	err := httpclient.GetJson(url, params, headers, &resp)
	if err != nil {
		logger.Error("%s request failed %v", url, err)
		return nil, err
	}
	//logger.Info("trc10 of url: %s and response: %+v", url, resp)
	var transfers []Trc10Data
	for _, transfer := range resp.Data {
		if transfer.ContractRet == "SUCCESS" && transfer.ContractData.Amount >= 10000 {
			if isCheckIn {
				if transfer.ToAddress == global.App.Config.Telegram.ReceiveAddress {
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
	if len(resp.Data) >= limit {
		var body Trc10Resp
		params["start"] = strconv.Itoa(page * limit)
		err = httpclient.GetJson(url, params, headers, &body)
		if err != nil {
			logger.Error("%s failed %v", url, err)
			return transfers, nil
		}
		page += 1
		if len(body.Data) >= 30 {
			for _, transfer := range body.Data {
				if transfer.ContractRet == "SUCCESS" && transfer.ContractData.Amount >= 10000 {
					if isCheckIn {
						if transfer.ToAddress == global.App.Config.Telegram.ReceiveAddress {
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
