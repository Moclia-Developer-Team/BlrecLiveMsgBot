package bot

import (
	"BliveBrocastPostBot/action"
	"BliveBrocastPostBot/data"
	"BliveBrocastPostBot/network"
	"errors"
	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"strconv"
	"strings"
)

// addRoom 添加主播指令实现
// 参数：pgid 群聊id/私聊id，user 触发操作的用户，mid 以room开头的roomid/以uid开头的uid，isprivate 是否是私聊
func addRoom(pgid int64, user int64, mid string, isPrivate bool) error {
	var rerr error = nil
	mark, err := data.CheckPrefix(mid, []string{"uid", "room"}) // 检查并读取mid前缀
	if err != nil {
		rerr = errors.New("输入的ID不合法" + mid)
		return rerr
	}
	mid = strings.TrimPrefix(mid, mark)
	id, _ := strconv.ParseInt(mid, 10, 64)
	// 根据mark分情况补全uid和roomid
	var roomid, userid int64
	switch mark {
	case "uid":
		userid = id
		// 根据UID获取RoomID
		roomid = network.GetRoomID(id)
		break
	case "room":
		roomid = id
		// 根据RoomID获取UID
		userid = network.GetUserID(id)
		break
	}
	// 根据UID和RoomID和pgid和私聊开关在配置文件中进行注册
	err = action.RegisterMember(pgid, user, userid, roomid, isPrivate)
	if err != nil {
		rerr = err
	}
	return rerr
}

func delRoom(pgid int64, mid string, isPrivate bool) error {
	var rerr error = nil
	mark, err := data.CheckPrefix(mid, []string{"uid", "room"}) // 检查并读取mid前缀
	if err != nil {
		rerr = errors.New("输入的ID不合法" + mid)
		return rerr
	}
	mid = strings.TrimPrefix(mid, mark)
	id, _ := strconv.ParseInt(mid, 10, 64)
	// 根据mark分情况补全uid和roomid
	var userid int64
	switch mark {
	case "room":
		// 根据RoomID获取UID
		userid = network.GetUserID(id)
		break
	}
	// 根据UID和RoomID和pgid和私聊开关在配置文件中进行注册
	err = action.DelRegisterConf(userid, pgid, isPrivate)
	if err != nil {
		rerr = err
	}
	return rerr
}

func atAllTriger(gid int64, mid string) error {
	var rerr error = nil
	mark, err := data.CheckPrefix(mid, []string{"uid", "room"}) // 检查并读取mid前缀
	if err != nil {
		rerr = errors.New("输入的ID不合法" + mid)
		return rerr
	}
	mid = strings.TrimPrefix(mid, mark)
	id, _ := strconv.ParseInt(mid, 10, 64)
	// 根据mark分情况补全uid和roomid
	var userid int64
	switch mark {
	case "room":
		// 根据RoomID获取UID
		userid = network.GetUserID(id)
		break
	}
	err = data.ChangeAtALL(gid, userid)
	if err != nil {
		rerr = err
	}
	return rerr
}

func listAllRooms() string {
	str := "所有已注册的主播列表：\n"
	rcd := <-data.RCData
	data.RCData <- rcd
	for _, user := range rcd.Users {
		str += strconv.FormatInt(user.UID, 10) + "(" + strconv.FormatInt(user.RoomID, 10) + "):\n"
		str += "Groups:\n"
		for _, group := range user.Groups {
			if group.AtAll {
				str += "\t" + strconv.FormatInt(group.Gid, 10) + ":开启At全体\n"
			} else {
				str += "\t" + strconv.FormatInt(group.Gid, 10) + ":关闭At全体\n"
			}
		}
	}
	return str
}

func addRoomToGroup(gid int64, superuser int64, mid string) error {
	err := addRoom(gid, superuser, mid, false)
	if err != nil {
		return err
	}
	return nil
}

func delRoomFromGroup(gid int64, mid string) error {
	err := delRoom(gid, mid, false)
	if err != nil {
		return err
	}
	return nil
}

//func groupAtAllTriger(gid int64) {
//
//}
//
//func recordTriger(mid string) {
//
//}

func defineAdminCommand() {
	zero.OnCommand("列出所有主播").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		if !(zero.OnlyPrivate(ctx) && zero.SuperUserPermission(ctx)) {
			return
		} // 非机器人管理组私聊无效
		msg := listAllRooms()
		ctx.Send(message.Message{message.Text(msg)})
	})

	//zero.OnCommand("列出群主播").SetBlock(true).Handle(func(ctx *zero.Ctx) {
	//	if !(zero.OnlyPrivate(ctx) && zero.SuperUserPermission(ctx)) {
	//		return
	//	} // 非机器人管理组私聊无效
	//
	//})

	//zero.OnCommand("添加主播到群").SetBlock(true).Handle(func(ctx *zero.Ctx) {
	//	if !(zero.OnlyPrivate(ctx) && zero.SuperUserPermission(ctx)) {
	//		return
	//	} // 非机器人管理组私聊无效
	//	dat := <-data.ConfData
	//	data.ConfData <- dat
	//	msgStr := ctx.MessageString()                                            // 接收完整消息内容
	//	msgStr = strings.TrimPrefix(msgStr, dat.BotConfig.CmdPrefix+"添加主播 ") // 清除前缀
	//	roomList := strings.Split(msgStr, " ")                                   // 将得到的参数以空格进行分割
	//})

	//zero.OnCommand("删除主播到群").SetBlock(true).Handle(func(ctx *zero.Ctx) {
	//	if !(zero.OnlyPrivate(ctx) && zero.SuperUserPermission(ctx)) {
	//		return
	//	} // 非机器人管理组私聊无效
	//
	//})

	//zero.OnCommand("开关群全体提醒").SetBlock(true).Handle(func(ctx *zero.Ctx) {
	//	if !(zero.OnlyPrivate(ctx) && zero.SuperUserPermission(ctx)) {
	//		return
	//	} // 非机器人管理组私聊无效
	//
	//})

	zero.OnCommand("开关主播录播").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		if !(zero.OnlyPrivate(ctx) && zero.SuperUserPermission(ctx)) {
			return
		} // 非机器人管理组私聊无效

	})
}

func DefineBotCommand() {
	log.Info("正在注册内置指令……")
	zero.OnCommand("添加主播").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		// 不是私聊或者管理员就不执行
		if !zero.UserOrGrpAdmin(ctx) {
			return
		}
		dat := <-data.ConfData
		data.ConfData <- dat
		msgStr := ctx.MessageString()                                        // 接收完整消息内容
		msgStr = strings.TrimPrefix(msgStr, dat.BotConfig.CmdPrefix+"添加主播 ") // 清除前缀
		roomList := strings.Split(msgStr, " ")                               // 将得到的参数以空格进行分割
		// 处理所有房间号/uid
		for _, adds := range roomList {
			if zero.OnlyGroup(ctx) {
				err := addRoom(ctx.Event.GroupID, ctx.Event.UserID, adds, false)
				if err != nil {
					ctx.Send(message.Message{message.Text("添加主播失败：" + adds)})
				} else {
					ctx.Send(message.Message{message.Text("添加主播成功：" + adds)})
				}
			} else if zero.OnlyPrivate(ctx) {
				err := addRoom(ctx.Event.UserID, ctx.Event.UserID, adds, true)
				if err != nil {
					ctx.Send(message.Message{message.Text("添加主播失败：" + adds)})
				} else {
					ctx.Send(message.Message{message.Text("添加主播成功：" + adds)})
				}
			} else {
				return
			}
		}
	})

	zero.OnCommand("删除主播").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		if !zero.UserOrGrpAdmin(ctx) {
			return
		}
		dat := <-data.ConfData
		data.ConfData <- dat
		msgStr := ctx.MessageString()                                        // 接收完整消息内容
		msgStr = strings.TrimPrefix(msgStr, dat.BotConfig.CmdPrefix+"删除主播 ") // 清除前缀
		roomList := strings.Split(msgStr, " ")                               // 将得到的参数以空格进行分割
		for _, dels := range roomList {
			if zero.OnlyGroup(ctx) {
				err := delRoom(ctx.Event.GroupID, dels, false)
				if err != nil {
					ctx.Send(message.Message{message.Text("删除主播失败：" + dels)})
				} else {
					ctx.Send(message.Message{message.Text("删除主播成功：" + dels)})
				}
			} else if zero.OnlyPrivate(ctx) {
				err := delRoom(ctx.Event.GroupID, dels, true)
				if err != nil {
					ctx.Send(message.Message{message.Text("删除主播失败：" + dels)})
				} else {
					ctx.Send(message.Message{message.Text("删除主播成功：" + dels)})
				}
			} else {
				return
			}
		}
	})

	zero.OnCommand("列出主播").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		if !zero.UserOrGrpAdmin(ctx) {
			return
		}
		if zero.OnlyGroup(ctx) {
			str := "本群已经注册的主播有：\n"
			for _, room := range data.GetRoomRegisters(ctx.Event.GroupID, false) {
				str += strconv.FormatInt(room, 10) + "\n"
			}
			ctx.Send(message.Message{message.Text(str)})
		} else if zero.OnlyPrivate(ctx) {
			str := "已经注册的主播有：\n"
			for _, room := range data.GetRoomRegisters(ctx.Event.UserID, true) {
				str += strconv.FormatInt(room, 10) + "\n"
			}
			ctx.Send(message.Message{message.Text(str)})
		}
	})

	zero.OnCommand("开关全体提醒").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		if !zero.AdminPermission(ctx) {
			return
		}
		dat := <-data.ConfData
		data.ConfData <- dat
		msgStr := ctx.MessageString()                                        // 接收完整消息内容
		msgStr = strings.TrimPrefix(msgStr, dat.BotConfig.CmdPrefix+"删除主播 ") // 清除前缀
		roomList := strings.Split(msgStr, " ")                               // 将得到的参数以空格进行分割
		for _, atall := range roomList {
			if zero.OnlyGroup(ctx) {
				err := atAllTriger(ctx.Event.GroupID, atall)
				if err != nil {
					ctx.Send(message.Message{message.Text("切换主播全体at提醒失败：" + atall)})
				} else {
					ctx.Send(message.Message{message.Text("切换主播全体at提醒成功：" + atall)})
				}
			} else {
				return
			}
		}
	})

	defineAdminCommand()
	log.Info("内置指令注册完成！")
}

// SendMessageByUser 通过主播UID确定将消息发往何处
// userid 主播的UID
// message 要发送的消息
func SendMessageByUser(userid int64, message string) {
	conf := <-data.ConfData
	data.ConfData <- conf
	gdata := data.GetUserGroupConfig(userid)
	pdata := data.GetUserPrivateConfig(userid)
	if len(gdata) == 0 && len(pdata) == 0 {
		return
	}
	for sendGroup, setAtAll := range gdata {
		if sendGroup == 0 {
			continue
		}
		var msg string
		if setAtAll {
			msg = "[CQ:At,all]" + message
		} else {
			msg = message
		}
		zero.GetBot(conf.BotConfig.Account).SendGroupMessage(sendGroup, msg)
	}

	for _, sendPrivate := range pdata {
		zero.GetBot(conf.BotConfig.Account).SendPrivateMessage(sendPrivate, message)
	}
}
