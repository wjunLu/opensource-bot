package group

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/eatmoreapple/openwechat"
	"github.com/samber/lo"
)

// 保存所有好友信息
var AllFriends = openwechat.Friends{}

// 保存每一位聊天对象上一次的有效command
var FriendsInChat = make(map[string]string)

// 记录可选的命令功能
var Commands = map[string]Command{
	"帮助": {
		Header: "您好，我是开源小助手!\n请回复以下关键词选择您需要的操作：\n",
		Options: Group{
			"加群": {
				Desc: "我要加入社区项目交流群",
			},
			"帮助": {
				Desc: "显示此条提示消息",
			},
		},
	},
	"加群": {
		Header:  "请回复以下群名称前的数字选择您需要加入的群组：\n",
		Options: make(Group),
		Handler: func(args ...interface{}) interface{} {
			return InviteFriendToGroup(args[0].(*openwechat.Self), args[1].(*openwechat.Friend), args[2].(*openwechat.Group))
		},
	},
}

func FileGroupsToCommands(groups openwechat.Groups) {
	cnt := 1
	for _, group := range groups {
		Commands["加群"].Options[strconv.Itoa(cnt)] = GroupInfo{
			Desc: group.NickName,
			Obj:  group,
		}
		cnt += 1
	}

}

func GetAllGroups(self *openwechat.Self) {
	groups, err := self.Groups()
	fmt.Println(groups, err)
	// 生成群组列表
	FileGroupsToCommands(groups)
}

func GetAllFriends(self *openwechat.Self) {
	friends, err := self.Friends()
	fmt.Println(friends, err)
	AllFriends = friends
}

func GetValidCommand(text string, commands map[string]Command) string {
	for key := range commands {
		// 检查字符串text是否包含map中的key
		if strings.Contains(text, key) {
			return key
		}
	}
	return ""
}

func GetReplyMessage(cmd Command) string {
	var text strings.Builder
	text.WriteString(cmd.Header)
	text.WriteString("\n")
	// 按照选项标号顺序生成Message字符串
	keys := lo.Keys(cmd.Options)
	sort.Strings(keys)
	for _, key := range keys {
		text.WriteString("【")
		text.WriteString(key)
		text.WriteString("】")
		text.WriteString(cmd.Options[key].Desc)
		text.WriteString("\n")
	}
	return text.String()
}

func InviteFriendToGroup(self *openwechat.Self, friend *openwechat.Friend, group *openwechat.Group) string {
	if group == nil {
		return "您所选的群组目标不存在！"
	}
	if err := self.AddFriendIntoManyGroups(friend, group); err != nil {
		fmt.Println(err)
	}
	group.SendText(fmt.Sprintf("@%s 欢迎入群！", friend.NickName))
	return "已成功邀请您加入目标群组！"
}
