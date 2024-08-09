package msg

import (
	"fmt"

	group "github.com/wjunlu/opensource-bot/group"
	"github.com/wjunlu/opensource-bot/record"

	"github.com/eatmoreapple/openwechat"
)

func HandleTextMessages(self *openwechat.Self, sender *openwechat.User, text string) string {
	key := group.GetValidCommand(text, group.Commands)
	// 根据收到信息判断当前操作是否为一级命令，进行不同处理
	if key == "" {
		// 非一级命令时，根据前置命令进行下一步操作
		if cmd, ok := group.FriendsInChat[sender.NickName]; ok {
			info, exists := group.Commands[cmd].Options[text]
			if exists {
				return group.Commands[cmd].Handler(self, &openwechat.Friend{User: sender}, info.Obj).(string)
			} else {
				return fmt.Sprintf("抱歉，所选目标不存在！\n请发送【%s】查看提示", cmd)
			}
		}
		return group.GetReplyMessage(group.Commands["帮助"])
	} else {
		// 将当前一级命令记录并给出操作提示
		group.FriendsInChat[sender.NickName] = key
		return group.GetReplyMessage(group.Commands[key])
	}
}

func HandleAllMessages(bot *openwechat.Bot, self *openwechat.Self) {
	bot.MessageHandler = func(msg *openwechat.Message) {
		reply := ""
		if msg.IsFriendAdd() {
			friend, err := msg.Agree()
			if err == nil {
				friend.SendText("您好，我是开源小助手，请发送【帮助】获取支持！")
				msg.AsRead()
				fmt.Printf("Succeed to add friend: %+v", friend)
				return
			}
			fmt.Println(err)
			return
		}
		if msg.IsComeFromGroup() {
			record.RecordGroupChat(msg)
			return
		}
		if !msg.IsText() {
			reply = "抱歉，请您发送文本消息！"
		} else {
			// 处理消息前总是更新群组列表
			group.GetAllGroups(self)
			// 获取消息发送者，处理后返回回复信息
			sender, _ := msg.Sender()
			reply = HandleTextMessages(self, sender, msg.Content)
			fmt.Println(msg.Content)
			fmt.Println(reply)
			msg.ReplyText(reply)
		}
	}
}
