package tron

import (
	"fmt"
	"telegram-trx/pkg/common/httpclient"
	"telegram-trx/pkg/core/cst"
	"telegram-trx/pkg/core/logger"
)

type AssetOverview struct {
	TotalAssetInTrx float64 `json:"totalAssetInTrx"`
	Data            []Asset `json:"data"`
	TotalTokenCount int     `json:"totalTokenCount"`
	TotalAssetInUsd float64 `json:"totalAssetInUsd"`
}

type Asset struct {
	TokenId         string  `json:"tokenId"`
	TokenName       string  `json:"tokenName"`
	TokenAbbr       string  `json:"tokenAbbr"`
	Vip             bool    `json:"vip"`
	Balance         string  `json:"balance"`
	TokenPriceInTrx float64 `json:"tokenPriceInTrx"`
	TokenPriceInUsd float64 `json:"tokenPriceInUsd"`
	AssetInTrx      float64 `json:"assetInTrx"`
	AssetInUsd      float64 `json:"assetInUsd"`
	Percent         float64 `json:"percent"`
}

func TokenAssetOverview(address string) (*Asset, error) {
	uri := fmt.Sprintf("%s%s", cst.TronScanApi, cst.TronWalletOverview)
	params := map[string]string{
		"address": address,
	}
	headers := map[string]string{}
	var resp AssetOverview
	err := httpclient.GetJson(uri, params, headers, &resp)
	if err != nil {
		logger.Error("[services] %s request failed %v", uri, err)
		return nil, err
	}
	logger.Info("[services] token asset overview: %+v", resp)
	for _, asset := range resp.Data {
		if asset.TokenName == "trx" && asset.TokenAbbr == "trx" {
			return &asset, nil
		}
	}
	return nil, nil
}
