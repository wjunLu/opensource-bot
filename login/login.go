package login

import (
	"fmt"

	"github.com/eatmoreapple/openwechat"
)

func WechatLogin() *openwechat.Bot {
	// 桌面模式
	bot := openwechat.DefaultBot(openwechat.Desktop)
	// 创建热存储容器对象
	reloadStorage := openwechat.NewFileHotReloadStorage("storage.json")
	// 热登录
	if err := bot.PushLogin(reloadStorage, openwechat.NewRetryLoginOption()); err != nil {
		fmt.Println(err)
		return nil
	}
	return bot
}

func GetLoginUser(bot *openwechat.Bot) *openwechat.Self {
	self, err := bot.GetCurrentUser()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return self
}
