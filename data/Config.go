package data

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
)

/***********************************************************************************************************************
=================================================配置文件和数据文件相关结构体=================================================
***********************************************************************************************************************/

type GroupData struct {
	Gid   int64 `json:"group_id"`
	AtAll bool  `json:"at_all"`
}

type RecordUsersData struct {
	UID      int64       `json:"uid"`
	RoomID   int64       `json:"room_id"`
	AddBy    int64       `json:"add_by"`
	Groups   []GroupData `json:"groups"`
	Privates []int64     `json:"privates"`
}

type RecordData struct {
	Users []RecordUsersData `json:"users"`
}

// ConfigData 配置文件结构体
type ConfigData struct {
	BotConfig struct {
		Name        string  `json:"name"`
		Account     int64   `json:"account"`
		Admin       []int64 `json:"admin"`
		URL         string  `json:"url"`
		AccessToken string  `json:"access_token"`
		CmdPrefix   string  `json:"command_prefix"`
	} `json:"bot_config"`
	BlrecConfig struct {
		WebhookAddr string `json:"webhook_addr"`
		WebhookPath string `json:"webhook_path"`
		ApiURL      string `json:"api_url"`
		APIKey      string `json:"api_key"`
	} `json:"blrec_config"`
}

var RCData = make(chan RecordData, 1)   // 全局数据文件
var ConfData = make(chan ConfigData, 1) // 全局配置文件
var RegData = make(chan []int64, 1)     // 已登录在配置中的主播列表

/***********************************************************************************************************************
=========================================================通用函数=========================================================
***********************************************************************************************************************/

// CheckFile 检查文件是否存在，-1无权限，0不存在，1存在，2是文件夹
func CheckFile(path string) int {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	if info.IsDir() {
		return 2
	}
	if info.Mode().Perm()&(1<<(uint(7))) == 0 {
		return -1
	}
	return 1
}

/***********************************************************************************************************************
====================================================配置文件读写相关函数====================================================
***********************************************************************************************************************/

// CreateConfigFile 创建默认的配置文件，创建成功返回0，创建失败返回1
func CreateConfigFile() int {
	var NewConfData = ConfigData{
		BotConfig: struct {
			Name        string  `json:"name"`
			Account     int64   `json:"account"`
			Admin       []int64 `json:"admin"`
			URL         string  `json:"url"`
			AccessToken string  `json:"access_token"`
			CmdPrefix   string  `json:"command_prefix"`
		}{
			Name:        "bot",
			Account:     0,
			Admin:       []int64{0},
			URL:         "ws://127.0.0.1:2280",
			AccessToken: "INITKEY1145141919810",
			CmdPrefix:   "/",
		},
		BlrecConfig: struct {
			WebhookAddr string `json:"webhook_addr"`
			WebhookPath string `json:"webhook_path"`
			ApiURL      string `json:"api_url"`
			APIKey      string `json:"api_key"`
		}{
			WebhookAddr: "127.0.0.1:2023",
			WebhookPath: "/post",
			ApiURL:      "http://127.0.0.1:2233",
			APIKey:      "bili2233"},
	}
	confFile, err := os.Create("config.json")
	if err != nil {
		return 1
	}
	configData, err := json.Marshal(NewConfData)
	if err != nil {
		return 1
	}
	var configStr bytes.Buffer
	err = json.Indent(&configStr, configData, "", "\t")
	if err != nil {
		return 0
	}
	_, err = confFile.WriteString(configStr.String())
	if err != nil {
		return 1
	}
	return 0
}

func ReadConfigFile() {
	var conf ConfigData
	file, err := os.ReadFile("config.json")
	if err != nil {
		return
	}
	err = json.Unmarshal(file, &conf)
	if err != nil {
		return
	}
	ConfData <- conf
}

func UpdateConfigFile() {
	updConfData, err := json.Marshal(ConfData)
	if err != nil {
		return
	}
	var confStr bytes.Buffer
	err = json.Indent(&confStr, updConfData, "", "\t")
	file, err := os.OpenFile("config.json", os.O_WRONLY|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	_, err = file.Seek(0, io.SeekEnd)
	if err != nil {
		return
	}
	_, err = file.WriteString(confStr.String())
	if err != nil {
		return
	}
}

/***********************************************************************************************************************
====================================================数据文件读写相关函数====================================================
***********************************************************************************************************************/

// CreateDataFile 创建默认的数据文件，创建成功返回0，创建失败返回1
func CreateDataFile() int {
	// 定义空白配置文件结构
	var NewData = RecordData{
		Users: []RecordUsersData{
			{
				UID:      0,
				RoomID:   0,
				AddBy:    0,
				Groups:   []GroupData{{Gid: 0, AtAll: false}},
				Privates: []int64{0},
			},
		},
	}
	// 创建空文件
	dataFile, err := os.Create("data.json")
	if err != nil {
		return 1
	}
	// 将数据文件结构解析成数组
	dataRepeat, err := json.Marshal(NewData)
	if err != nil {
		return 1
	}
	var dataStr bytes.Buffer
	// 格式化数据文件
	err = json.Indent(&dataStr, dataRepeat, "", "\t")
	// 将数据文件内容写入磁盘
	_, err = dataFile.WriteString(dataStr.String())
	if err != nil {
		return 1
	}
	return 0
}

// ReadDataFile 读取配置文件内容写入通道
func ReadDataFile() {
	var rcd RecordData
	file, err := os.ReadFile("data.json")
	if err != nil {
		return
	}
	err = json.Unmarshal(file, &rcd)
	if err != nil {
		return
	}
	RCData <- rcd
}

// UpdateDataFile 将通道内的缓存数据读取出来并写入配置文件
func UpdateDataFile() {
	rcd := <-RCData
	RCData <- rcd
	updRCData, err := json.Marshal(rcd)
	if err != nil {
		return
	}
	var dataStr bytes.Buffer
	err = json.Indent(&dataStr, updRCData, "", "\t")
	file, err := os.OpenFile("data.json", os.O_WRONLY|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	_, err = file.Seek(0, io.SeekEnd)
	if err != nil {
		return
	}
	_, err = file.WriteString(dataStr.String())
	if err != nil {
		return
	}
}

// AddNewChatToMember 将新群/私聊添加到主播的配置列表
// uid 主播的用户id，mid 私聊/群聊的号码，private 是否是私聊
func AddNewChatToMember(uidPos int, mid int64, private bool) {
	rcd := <-RCData
	if private {
		rcd.Users[uidPos].Privates = append(rcd.Users[uidPos].Privates, mid)
	} else {
		rcd.Users[uidPos].Groups = append(rcd.Users[uidPos].Groups, GroupData{
			Gid:   mid,
			AtAll: false,
		})
	}
	RCData <- rcd
	UpdateDataFile()
}

/***********************************************************************************************************************
=====================================================信息查询相关函数======================================================
***********************************************************************************************************************/

// GetUserGroupConfig 获取主播UID下的群组信息
func GetUserGroupConfig(uid int64) map[int64]bool {
	rcd := <-RCData
	RCData <- rcd
	grpupsData := make(map[int64]bool)
	for _, users := range rcd.Users {
		if users.UID == uid {
			for _, groups := range users.Groups {
				grpupsData[groups.Gid] = groups.AtAll
			}
			break
		}
	}
	return grpupsData
}

// GetUserPrivateConfig 获取主播UID下的私聊注册信息
func GetUserPrivateConfig(uid int64) []int64 {
	rcd := <-RCData
	RCData <- rcd
	var privateData []int64
	for _, users := range rcd.Users {
		if users.UID == uid {
			for _, private := range users.Privates {
				privateData = append(privateData, private)
			}
			break
		}
	}
	return privateData
}

// GetRegisterList 获取注册的主播列表
func GetRegisterList() {
	rcd := <-RCData
	RCData <- rcd
	var RegistedData []int64
	for _, users := range rcd.Users {
		RegistedData = append(RegistedData, users.UID)
	}
	select {
	case _ = <-RegData:
		RegData <- RegistedData
	default:
		RegData <- RegistedData
	}
}

// GetRoomRegisters 获取私聊/群聊注册的主播信息
func GetRoomRegisters(mid int64, isPrivate bool) []int64 {
	rcd := <-RCData
	RCData <- rcd
	var userId []int64
	for _, RegisterUid := range rcd.Users {
		if isPrivate {
			for _, priv := range RegisterUid.Privates {
				if priv == mid {
					userId = append(userId, RegisterUid.UID)
					break
				}
			}
		} else {
			for _, UidGroup := range RegisterUid.Groups {
				if UidGroup.Gid == mid {
					userId = append(userId, RegisterUid.UID)
					break
				}
			}
		}
	}
	return userId
}

// IsRegisted 判断主播是否已经在配置文件中注册
// 如果已注册，返回他在数组中的位置，未注册返回-1
func IsRegisted(uid int64) int {
	regd := <-RegData
	RegData <- regd
	for pos, uids := range regd {
		if uids == uid {
			return pos
		}
	}
	return -1
}

// ChangeAtALL 改变主播配置下特定群的at全体状态
func ChangeAtALL(gid int64, uid int64) error {
	uidPos := IsRegisted(uid)
	if uidPos == -1 {
		return errors.New("该用户未注册")
	}
	rcd := <-RCData
	for _, group := range rcd.Users[uidPos].Groups {
		if group.Gid == gid {
			group.AtAll = !group.AtAll
		}
	}
	RCData <- rcd
	UpdateDataFile()
	return nil
}
