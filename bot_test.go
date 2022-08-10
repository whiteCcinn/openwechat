package openwechat

import (
	"fmt"
	"testing"
	"time"
)

func TestHotLogin(t *testing.T) {
	bot := DefaultBot(DesktopInTerminal)

	dispatcher := NewMessageMatchDispatcher()
	dispatcher.SetAsync(true)

	// 只处理消息类型为文本类型的消息
	dispatcher.OnText(func(msg *MessageContext) {
		fmt.Println(msg.Message.FromUserName, msg.Message.ToUserName, msg.Message.Content)
		fmt.Println(msg.Message.RowContent)
		fmt.Println(msg.Message.OriContent)

		fmt.Println(msg.Message.IsComeFromGroup())
		fmt.Println(msg.Message.IsJoinGroup())
		fmt.Println(msg.Message.IsPaiYiPai())
		//msg.ReplyText("hello")
	})


	// 注册消息处理函数
	bot.MessageHandler = DispatchMessage(dispatcher)

	// 注册登陆二维码回调
	//UUIDCallback := func(model Mode, uuid string) {
	//	if !bot.IsHot() {
	//		PrintlnQrcodeUrl(model, uuid)
	//	}
	//}
	bot.UUIDCallback = PrintlnQrcodeUrl


	Storage := NewJsonFileHotReloadStorage("Storage.json")

	// 登陆
	if err := bot.HotLogin(Storage, true); err != nil {
		fmt.Println(err)
		return
	}

	// 获取登陆的用户
	self, err := bot.GetCurrentUser()
	if err != nil {
		fmt.Println(err)
		return
	}



	// 获取所有的好友
	friends, err := self.Friends()
	fmt.Println(friends, err)

	self.Detail()

	// 获取所有的群组
	groups, err := self.Groups(true)
	fmt.Println(groups, err)

	//info := bot.Storage.LoginInfo
	//members, err := bot.Caller.WebWxBatchGetContactGroup(info)
	//fmt.Println(members)

	// 阻塞主goroutine, 直到发生异常或者用户主动退出
	bot.Block()
}

func TestLogin(t *testing.T) {
	bot := DefaultBot(Desktop)
	bot.LoginCallBack = func(body []byte) {
		t.Log("login success")
	}
	if err := bot.Login(); err != nil {
		t.Error(err)
	}
}

func TestLogout(t *testing.T) {
	bot := DefaultBot(Desktop)
	bot.LoginCallBack = func(body []byte) {
		t.Log("login success")
	}
	bot.LogoutCallBack = func(bot *Bot) {
		t.Log("logout")
	}
	bot.MessageHandler = func(msg *Message) {
		if msg.IsText() && msg.Content == "logout" {
			bot.Logout()
		}
	}
	if err := bot.Login(); err != nil {
		t.Error(err)
		return
	}
	bot.Block()
}

func TestMessageHandle(t *testing.T) {
	bot := DefaultBot(Desktop)
	bot.MessageHandler = func(msg *Message) {
		if msg.IsText() && msg.Content == "ping" {
			msg.ReplyText("pong")
		}
	}
	if err := bot.Login(); err != nil {
		t.Error(err)
		return
	}
	bot.Block()
}

func TestFriends(t *testing.T) {
	bot := DefaultBot(Desktop)
	if err := bot.Login(); err != nil {
		t.Error(err)
		return
	}
	user, err := bot.GetCurrentUser()
	if err != nil {
		t.Error(err)
		return
	}
	friends, err := user.Friends()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(friends)
}

func TestGroups(t *testing.T) {
	bot := DefaultBot(Desktop)
	if err := bot.Login(); err != nil {
		t.Error(err)
		return
	}
	user, err := bot.GetCurrentUser()
	if err != nil {
		t.Error(err)
		return
	}
	groups, err := user.Groups()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(groups)
}

func TestPinUser(t *testing.T) {
	bot := DefaultBot(Desktop)
	if err := bot.Login(); err != nil {
		t.Error(err)
		return
	}
	user, err := bot.GetCurrentUser()
	if err != nil {
		t.Error(err)
		return
	}
	friends, err := user.Friends()
	if err != nil {
		t.Error(err)
		return
	}
	if friends.Count() > 0 {
		f := friends.First()
		f.Pin()
		time.Sleep(time.Second * 5)
		f.UnPin()
	}
}

func TestSender(t *testing.T) {
	bot := DefaultBot(Desktop)
	bot.MessageHandler = func(msg *Message) {
		if msg.IsSendByGroup() {
			fmt.Println(msg.SenderInGroup())
		} else {
			fmt.Println(msg.Sender())
		}
	}
	if err := bot.Login(); err != nil {
		t.Error(err)
		return
	}
	bot.Block()
}


