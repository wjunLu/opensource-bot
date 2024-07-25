package record

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/eatmoreapple/openwechat"
)

const ROOT_DOR = "/tmp/ascend-bot/record/groups/"

func RecordPictureMsg(msg *openwechat.Message) string {
	if !msg.IsPicture() {
		return msg.MsgId
	}
	// 创建文件
	file, err := os.Create(path.Join(ROOT_DOR, "images/"+msg.MsgId+".jpg"))
	if err != nil {
		fmt.Println(err)
		return msg.MsgId
	}
	defer file.Close()
	// 将响应体写入文件
	resp, err := msg.GetPicture()
	if err != nil {
		fmt.Println(err)
		return msg.MsgId
	}
	defer resp.Body.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println(err)
		return file.Name()
	}
	return msg.MsgId
}

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
	now := time.Now()
	timeStr := now.Format("2006-01-02 15:04:05")
	timeDayStr := now.Format("2006-01-02")
	fmt.Println(timeStr, timeDayStr)
	if !msg.IsText() {
		content = RecordPictureMsg(msg)
	}
	conversation := GroupMsgRecord{
		Time:    timeStr,
		Sender:  sender.NickName,
		Content: content,
	}
	// 当前消息存入文件，文件名为群名
	groupPath := path.Join(ROOT_DOR, group.NickName)
	if err := os.MkdirAll(groupPath, 0755); err != nil {
		fmt.Println(err)
		return
	}
	WriteDataIntoFile(&conversation, path.Join(groupPath, timeDayStr+".json"))
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
