package global

import (
	"github.com/go-redis/redis/v8"
	"github.com/mr-linch/go-tg"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"telegram-trx/pkg/config"
)

type Application struct {
	Config *config.Configuration
	DB     *gorm.DB
	Redis  *redis.Client
	Cron   *cron.Cron
	Client *tg.Client
}

var App = new(Application)
