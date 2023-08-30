package bootstrap

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"telegram-trx/pkg/core/global"
)

func LoadConfig(appName string) *viper.Viper {
	// 初始化 viper
	v := viper.New()
	v.SetConfigName(appName)
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("read config file %s.yaml failed: %s \n", appName, err))
	}

	// 监听配置文件
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		fmt.Printf("config file changed %s\n", appName)
		// 重载配置
		if err := v.Unmarshal(&global.App.Config); err != nil {
			fmt.Println(err)
		}
	})
	// 将配置赋值给全局变量
	if err := v.Unmarshal(&global.App.Config); err != nil {
		panic(err)
	}
	return v
}
