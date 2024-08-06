package main

import (
	"fmt"
	"os"

	"github.com/eatmoreapple/openwechat"
	"github.com/wjunlu/ascend-bot/record"

	"github.com/robfig/cron/v3"
	"github.com/wjunlu/ascend-bot/login"
	msg "github.com/wjunlu/ascend-bot/message"

	group "github.com/wjunlu/ascend-bot/group"
)

func main() {
	// 创建聊天记录存储目录
	err := os.MkdirAll(record.ROOT_DOR, 0755)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 登录微信账户
	bot := login.WechatLogin()
	if bot == nil {
		return
	}
	// 获取登陆的用户
	self := login.GetLoginUser(bot)
	// 注册消息处理函数
	msg.HandleAllMessages(bot, self)
	// 获取所有的好友
	group.GetAllFriends(self)
	// 启动定时任务
	RunPeriodicTasks()
	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}

func RunPeriodicTasks() {
	c := cron.New()
	c.AddFunc("@every 1h", func() {
		// 定时清空每天的聊天缓存
		group.FriendsInChat = map[*openwechat.User]string{}
	})
	c.Start()
}
