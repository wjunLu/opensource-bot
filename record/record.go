package record

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/eatmoreapple/openwechat"
)

const ROOT_DOR = "/tmp/ascend-bot/record/groups/"

func RecordGroupChat(msg *openwechat.Message) {
	user, err := msg.Sender()
	if err != nil {
		fmt.Println(err)
		return
	}
	group := &openwechat.Group{User: user}
	sender, err := msg.SenderInGroup()
	if err != nil {
		fmt.Println(err)
		return
	}
	content := msg.Content
	if !msg.IsText() {
		content = "|> Non-Text <|"
	}
	conversation := GroupMsgRecord{
		Sender:  sender.NickName,
		Content: content,
	}
	// 当前消息存入文件，文件名为群名
	WriteDataIntoFile(&conversation, path.Join(ROOT_DOR, group.NickName+".json"))
}

func WriteDataIntoFile(conversation *GroupMsgRecord, file string) {
	// 尝试以追加模式打开文件
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	// 将JSON数据编码为字节数组
	jsonData, err := json.Marshal(*conversation)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 将JSON数据追加到文件中
	_, err = f.Write(jsonData)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 换行符,以便更好地格式化JSON数据
	_, err = fmt.Fprintln(f)
	if err != nil {
		fmt.Println(err)
		return
	}
}
