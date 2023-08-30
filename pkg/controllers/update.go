package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap/buffer"
	"html/template"
	"net/http"
	"sync"
	"telegram-trx/pkg/core/cst"
	"telegram-trx/pkg/core/global"
	"telegram-trx/pkg/core/logger"
	"telegram-trx/pkg/models"
)

func Update(token string) {
	var users []models.User
	global.App.DB.Find(&users)
	tmpl, _ := template.ParseFiles(cst.UpdateTemplateFile)
	buf := new(buffer.Buffer)
	_ = tmpl.Execute(buf, nil)
	var wg sync.WaitGroup
	for _, user := range users {
		wg.Add(1)
		go func(user models.User) {
			defer wg.Done()
			uri := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)
			body := map[string]string{
				"chat_id":    user.UserID,
				"text":       buf.String(),
				"parse_mode": "HTML",
			}
			b, _ := json.Marshal(&body)
			resp, err := http.Post(uri, "application/json", bytes.NewReader(b))
			if err != nil {
				logger.Error("[update] send message to user failed %v", err)
			}
			defer func() { _ = resp.Body.Close() }()
		}(user)
	}
	wg.Wait()

}
