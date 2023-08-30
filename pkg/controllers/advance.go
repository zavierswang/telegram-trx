package controllers

import (
	"context"
	"github.com/mr-linch/go-tg"
	"github.com/mr-linch/go-tg/tgb"
	"go.uber.org/zap/buffer"
	"html/template"
	"strings"
	"telegram-trx/pkg/core/cst"
	"telegram-trx/pkg/core/logger"
	"telegram-trx/pkg/middleware"
)

func Advance(ctx context.Context, update *tgb.MessageUpdate) error {
	userId := update.Message.From.ID.PeerID()
	username := update.Message.From.Username.PeerID()
	username = strings.ReplaceAll(username, "@", "")
	logger.Info("[%s %s] trigger action [advance] controller", userId, username)
	middleware.SessionManager.Get(ctx).Step = middleware.SessionAdvance
	tmpl, err := template.ParseFiles(cst.AdvanceTemplateFile)
	if err != nil {
		logger.Error("[%s %s] template parse file %s, failed %v", userId, username, err)
		return err
	}
	buf := new(buffer.Buffer)
	err = tmpl.Execute(buf, nil)
	if err != nil {
		logger.Error("[%s %s] template execute file %s, failed %v", userId, username, err)
		return err
	}
	return update.Update.
		Reply(ctx, update.Answer(buf.String()).
			ParseMode(tg.HTML).
			DisableWebPagePreview(true))
}
