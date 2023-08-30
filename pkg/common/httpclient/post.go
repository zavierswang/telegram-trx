package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap/buffer"
	"io"
	"net/http"
	"net/url"
	"telegram-trx/pkg/core/logger"
	"time"
)

func PostJson(uri string, body interface{}, headers map[string]string, params map[string]string, data interface{}) error {
	values := url.Values{}
	for key, value := range params {
		values.Set(key, value)
	}
	uri = fmt.Sprintf("%s?%s", uri, values.Encode())
	//logger.Info("PostJson: %s", uri)
	buf, _ := json.Marshal(&body)
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewReader(buf))
	if err != nil {
		return err
	}
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	client := http.Client{
		Transport:     &http.Transport{},
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       time.Second * 60,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		logger.Error("%s request errCode: %d", uri, resp.StatusCode)
		return err
	}
	b := new(buffer.Buffer)
	if _, err := io.Copy(b, resp.Body); err != nil {
		logger.Error("%s writer buff failed %v", err)
		return err
	}
	//logger.Info("%s response successfully, response: %s", uri, b.String())
	return json.Unmarshal(b.Bytes(), &data)
}
