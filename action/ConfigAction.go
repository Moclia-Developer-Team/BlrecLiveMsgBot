package action

import (
	"BliveBrocastPostBot/data"
	"BliveBrocastPostBot/network"
	"errors"
)

// DelRegisterConf 将群或私聊从主播配置下移除，若主播没被监听则完全删除
func DelRegisterConf(uid int64, pgid int64, isPrivate bool) error {
	rcd := <-data.RCData // 将通道内的数据拉下来
	pos := data.IsRegisted(uid)
	isDeleted := false
	if pos == -1 {
		return errors.New("主播配置文件不存在")
	}
	// 查找并删除群或私聊在主播配置下的位置
	if isPrivate {
		for index, pid := range rcd.Users[pos].Privates { // 查找要删除的私聊频道的位置并删除
			if pid == pgid {
				rcd.Users[pos].Privates =
					append(rcd.Users[pos].Privates[:index], rcd.Users[pos].Privates[index+1:]...) // 将号码从私聊中删除
				isDeleted = true
				break
			}
		}
	} else {
		for index, gid := range rcd.Users[pos].Groups { // 查找要删除的群聊频道的位置并删除
			if gid.Gid == pgid {
				rcd.Users[pos].Groups =
					append(rcd.Users[pos].Groups[:index], rcd.Users[pos].Groups[index+1:]...) // 将号码从私聊中删除
				isDeleted = true
				break
			}
		}
	}
	if !isDeleted {
		data.RCData <- rcd
		return errors.New("主播未在本频道中注册")
	}
	// 检查配置发起者，如果不是管理员发起就同时删除blrec记录
	if rcd.Users[pos].AddBy != 0 {
		network.BlrecDelRoom(rcd.Users[pos].RoomID)
	}
	// 检查是否存在其他监听地址，不存在就删除配置
	if len(rcd.Users[pos].Privates) == 0 && len(rcd.Users[pos].Groups) == 1 {
		rcd.Users = append(rcd.Users[:pos], rcd.Users[pos+1:]...) // 将主播配置从配置文件中删除
	}
	data.RCData <- rcd
	// 保存新配置文件，生成新的主播列表
	data.UpdateDataFile()
	data.GetRegisterList()
	return nil
}

// RegisterNewMember 将新主播添加到配置文件
func RegisterNewMember(adder int64, uid int64, roomid int64) {
	rcd := <-data.RCData
	// 判断录播姬中是否存在用户
	if network.BlrecCheckRoomAvailable(roomid) { // 如果存在adder就为0
		adder = 0
	} else { // 不存在就在录播姬中添加
		network.BlrecAddRoom(roomid)
		network.BlrecStopAutoRecord(roomid)
	}
	// 在配置文件中添加用户的默认配置
	rcd.Users = append(rcd.Users, data.RecordUsersData{
		UID:    uid,
		RoomID: roomid,
		AddBy:  adder,
		Groups: []data.GroupData{
			{
				Gid:   0,
				AtAll: false,
			},
		},
		Privates: []int64{},
	})
	data.RCData <- rcd // 将修改后的配置文件更新到通道
	// 保存新配置文件，生成新的主播列表
	data.UpdateDataFile()
	data.GetRegisterList()
}

// RegisterMember 更新主播的配置文件/新建主播配置
func RegisterMember(mid int64, user int64, uid int64, roomid int64, isPrivate bool) error {
	uidPos := data.IsRegisted(uid) // 判断主播是否已经注册
	if uidPos != -1 {              // 已经注册直接添加配置
		if data.CheckUsersIsRegisted(mid, uidPos, isPrivate) {
			return errors.New("主播已在本频道注册")
		}
		data.AddNewChatToMember(uidPos, mid, isPrivate)
	} else { // 未被注册先注册后添加
		RegisterNewMember(user, uid, roomid)
		uidPos = data.IsRegisted(uid)
		data.AddNewChatToMember(uidPos, mid, isPrivate)
	}
	return nil
}
