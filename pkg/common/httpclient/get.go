package httpclient

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap/buffer"
	"io"
	"net/http"
	"net/url"
	"telegram-trx/pkg/core/logger"
	"time"
)

func GetJson(uri string, params map[string]string, headers map[string]string, data interface{}) error {
	values := url.Values{}
	for key, value := range params {
		values.Set(key, value)
	}
	uri = fmt.Sprintf("%s?%s", uri, values.Encode())
	//logger.Info("GetJson: %s", uri)
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		logger.Error("%s req failed %w", uri, err)
		return err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	//logger.Info("%s", uri)
	httpClient := &http.Client{
		Timeout:   30 * time.Second,
		Transport: &http.Transport{},
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Error("%s request failed %w", uri, err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		logger.Error("%s request errCode: %d", uri, resp.StatusCode)
		return err
	}
	buf := new(buffer.Buffer)
	if _, err := io.Copy(buf, resp.Body); err != nil {
		logger.Error("%s writer buff failed %v", err)
		return err
	}
	//logger.Info("%s response successfully", uri)
	return json.Unmarshal(buf.Bytes(), &data)
}
