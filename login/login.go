package login

import (
	"fmt"
	"runtime"

	"github.com/eatmoreapple/openwechat"
	"github.com/skip2/go-qrcode"
)

func WechatLogin() *openwechat.Bot {
	// 桌面模式
	bot := openwechat.DefaultBot(openwechat.Desktop)
	// linux终端打印二维码
	if runtime.GOOS == "linux" {
		bot.UUIDCallback = ConsoleQrCode
	}
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

func ConsoleQrCode(uuid string) {
	q, _ := qrcode.New("https://login.weixin.qq.com/l/"+uuid, qrcode.Low)
	fmt.Println(q.ToString(true))
}
