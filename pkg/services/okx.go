package services

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap/buffer"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"telegram-trx/pkg/core/cst"
	"telegram-trx/pkg/core/logger"
	"time"
)

type RateResp struct {
	Code string `json:"code"`
	Data []Rate `json:"data"`
}

type Rate struct {
	InstId string `json:"instId"`
	Side   string `json:"side"`
	Px     string `json:"px"`
	Ts     string `json:"ts"`
}

func GetOkxMarketRate() (*RateResp, error) {
	params := map[string]string{
		"t":      strconv.FormatInt(time.Now().UnixMilli(), 10),
		"instId": "TRX-USDT",
	}
	headers := map[string]string{
		"User-Agent":             cst.UserAgent,
		"Content-Type":           "application/json;charset=UTF-8",
		"Transfer-Encoding":      "chunked",
		"Vary":                   "Accept-Encoding",
		"X-Content-Type-Options": "nosniff",
		"Cache-Control":          "no-cache",
		"Content-Encoding":       "gzip",
		"Cookie":                 "preferLocale=zh_CN; devId=0f7e8aab-7acc-4f5a-8de4-fb4df4031f4a; first_ref=https%3A%2F%2Fwww.google.com.hk%2F; intercom-id-ny9cf50h=6911d9a0-594e-44fe-8d59-d2b22401001e; intercom-device-id-ny9cf50h=733e1c57-c180-4eb2-9c15-f65de162b306; G_ENABLED_IDPS=google; isLogin=1; _ga=GA1.1.367052617.1684696743; _ga_G0EKWWQGTZ=GS1.1.1684747969.3.0.1684747969.60.0.0; locale=zh_CN; okg.currentMedia=xl; token=eyJhbGciOiJIUzUxMiJ9.eyJqdGkiOiJleDExMDE2ODQ0ODY3MDQ5NDVFNzYyRjVBMUM3OEVCMDU4MXBXR0oiLCJ1aWQiOiIvK2JEYzdFdXFZZ0FvZVVZSkZXWDZ3PT0iLCJzdGEiOjAsIm1pZCI6Ii8rYkRjN0V1cVlnQW9lVVlKRldYNnc9PSIsImlhdCI6MTY4NTAyNjExMiwiZXhwIjoxNjg1NjMwOTEyLCJiaWQiOjAsImRvbSI6Ind3dy5va3guY29tIiwiZWlkIjoxLCJpc3MiOiJva2NvaW4iLCJzdWIiOiJENkZEMkRBRjM5MUQ3Q0Q3RkJBMEQ4MEZGMjk2MThFNCJ9.ijE6nl9TVPfHkhsbvkIBsllOvP6QvDhJJASJZA38HioHWi2rg_mKNcLwmd6CvhjmjOb5unR6u1JYOSVbF1ZBFw; __cf_bm=GLe7ax74hv0kV7ISk1J42jvlDJnXKso5yQ1Zweui8BY-1685026594-0-AUYno84bVV0D2OO64n8SlAq7A46zdJyLdj308VhJgHN5s5w/rfZw+e4qTUK44giI5p7kpdkAFwCl6MKkkR4dr4M=; intercom-session-ny9cf50h=VE9DR09IbG9ITFd2bmtjS1pCLzVwZ2huOGJEWTNoZ3JhVlVQVmNrWFJiZGgxU1BWRHljUlV6SXZQaXoyNTNwVS0tdEpXeXpLNU8rWGVyOTNjMkpqWTlxZz09--a20ae0214384a3cc4f0854977f9ad3c030c2935f; ok-ses-id=aNl+fd0J3L4WCLX4SZoLXkpNke4sJL4Zmp6VAVhNDL+fgtydnOZJDtHL1MXXpanWGg5o/6WY4q/o+WNrdVUMVJC1CqCqSAaufi4EXmM2yKJ+Ws5SYWQ33vcxSh5zoyh/; _monitor_extras={\"deviceId\":\"2BLFQ1_STxGSDC_Ml2gGlU\",\"eventId\":129,\"sequenceNumber\":129}; amp_56bf9d=9AgAMDEwwHb6LqO6dogW9w.ZHREN0xpeDlxZ21jL2pSLzZ3UzM4Zz09..1h19m467s.1h19n3d0h.2e.8.2m",
		"Referer":                "https://www.okx.com/cn/trade-spot/trx-usdt",
		"Accept":                 "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
		"Authorization":          "eyJhbGciOiJIUzUxMiJ9.eyJqdGkiOiJleDExMDE2ODQ0ODY3MDQ5NDVFNzYyRjVBMUM3OEVCMDU4MXBXR0oiLCJ1aWQiOiIvK2JEYzdFdXFZZ0FvZVVZSkZXWDZ3PT0iLCJzdGEiOjAsIm1pZCI6Ii8rYkRjN0V1cVlnQW9lVVlKRldYNnc9PSIsImlhdCI6MTY4NTAyNjExMiwiZXhwIjoxNjg1NjMwOTEyLCJiaWQiOjAsImRvbSI6Ind3dy5va3guY29tIiwiZWlkIjoxLCJpc3MiOiJva2NvaW4iLCJzdWIiOiJENkZEMkRBRjM5MUQ3Q0Q3RkJBMEQ4MEZGMjk2MThFNCJ9.ijE6nl9TVPfHkhsbvkIBsllOvP6QvDhJJASJZA38HioHWi2rg_mKNcLwmd6CvhjmjOb5unR6u1JYOSVbF1ZBFw",
	}
	var body RateResp

	values := url.Values{}
	for key, value := range params {
		values.Set(key, value)
	}
	uri := fmt.Sprintf("%s?%s", cst.OkxMarketTradesApi, values.Encode())

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		logger.Error("%s request failed %w", uri, err)
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	client := http.Client{Timeout: time.Second * 30}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("%s http.Do failed %v", uri, err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		logger.Error("%s http.Do errCode: %d", uri, resp.StatusCode)
		return nil, err
	}
	buf := new(buffer.Buffer)
	if _, err := io.Copy(buf, resp.Body); err != nil {
		logger.Error("%s read body failed %w", uri, err)
		return nil, err
	}
	if err := json.Unmarshal(buf.Bytes(), &body); err != nil {
		logger.Error("%s unmarshal failed %w", uri, err)
		return nil, err
	}
	return &body, nil
}

type OKXResp struct {
	Code int      `json:"code"`
	Data RateDate `json:"data"`
}

type RateDate struct {
	Sell []RateSell `json:"sell"`
}

type RateSell struct {
	NickName     string `json:"nickName"`
	Price        string `json:"price"`
	QuoteSymbol  string `json:"quoteSymbol"`
	BaseCurrency string `json:"baseCurrency"`
}

func GetOkxTradingOrders() (*OKXResp, error) {
	params := map[string]string{
		"t":                 strconv.FormatInt(time.Now().UnixMilli(), 10),
		"quoteCurrency":     "CNY",
		"baseCurrency":      "USDT",
		"side":              "sell",
		"paymentMethod":     "all",
		"userType":          "all",
		"showTrade":         "false",
		"receivingAds":      "false",
		"showFollow":        "false",
		"showAlreadyTraded": "false",
		"isAbleFilter":      "false",
	}
	headers := map[string]string{
		"User-Agent":                cst.UserAgent,
		"Content-Type":              "application/json;charset=UTF-8",
		"Transfer-Encoding":         "chunked",
		"Vary":                      "Accept-Encoding",
		"X-Content-Type-Options":    "nosniff",
		"Cache-Control":             "no-cache, no-store, max-age=0, must-revalidate",
		"X-Frame-Options":           "DENY",
		"Strict-Transport-Security": "max-age=63072000; includeSubdomains; preload",
		"Content-Encoding":          "gzip",
	}
	var body OKXResp

	values := url.Values{}
	for key, value := range params {
		values.Set(key, value)
	}
	uri := fmt.Sprintf("%s?%s", cst.OkxTradingOrdersApi, values.Encode())

	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		logger.Error("%s request failed %w", uri, err)
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	client := http.Client{Timeout: time.Second * 30}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("%s http.Do failed %v", uri, err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		logger.Error("%s http.Do errCode: %d", uri, resp.StatusCode)
		return nil, err
	}
	buf := new(buffer.Buffer)
	if _, err := io.Copy(buf, resp.Body); err != nil {
		logger.Error("%s read body failed %w", uri, err)
		return nil, err
	}
	if err := json.Unmarshal(buf.Bytes(), &body); err != nil {
		logger.Error("%s unmarshal failed %w", uri, err)
		return nil, err
	}
	return &body, nil
}
