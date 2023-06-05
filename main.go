package main

import (
	"BliveBrocastPostBot/bot"
	"BliveBrocastPostBot/data"
	"BliveBrocastPostBot/network"
	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
	"os"
	"strconv"
)

/*
function：
检查配置文件是否存在
检查数据文件是否存在
存在配置文件读取配置文件
存在数据文件读取数据文件
不存在任意文件自动创建空白文件
*/
func init() {
	switch data.CheckFile("data/config.json") {
	case -1:
		// 无法操作，直接退出
		log.Fatal("无目录操作权限！")
		break
	case 0:
		// 创建配置文件内容
		log.Warn("未检测到配置文件，正在创建空配置文件……")
		data.CreateConfigFile()
		log.Warn("配置文件创建完成！请修改后重新启动")
		os.Exit(0)
	case 1:
		// 读取配置文件内容
		data.ReadConfigFile()
		log.Info("配置文件读取完成！")
		break
	case 2:
		log.Fatal("存在同名文件夹！")
		break
	}

	switch data.CheckFile("data/data.json") {
	case -1:
		// 无法操作，直接退出
		log.Fatal("无目录操作权限！")
		break
	case 0:
		// 创建配置文件内容，并读取
		log.Warn("未检测到数据文件，正在创建空数据文件……")
		data.CreateDataFile()
		log.Warn("数据文件创建完成！")
		data.ReadDataFile()
		break
	case 1:
		// 读取数据文件内容
		data.ReadDataFile()
		log.Info("数据文件读取完成！")
		break
	case 2:
		log.Fatal("存在同名文件夹！")
		break
	}
	log.Info("正在读取主播列表……")
	data.GetRegisterList() // 更新登录的主播列表
	reg := <-data.RegData
	data.RegData <- reg
	for _, list := range reg {
		log.Info("已读取：" + strconv.FormatInt(list, 10))
	}
	log.Info("主播列表读取完成")
	network.UpdateMixinKeyOnStartup()
	log.Info("初始化检查完成！")
}

func main() {
	conf := <-data.ConfData
	data.ConfData <- conf // 读完后立刻还回去给其他分支用
	log.Info("开始运行主程序")
	//log.Info(conf)
	go network.UpdateMixinKey()
	// 启动http服务器
	go network.StartWeb(conf.BlrecConfig.WebhookAddr, conf.BlrecConfig.WebhookPath)
	// 启动zerobot
	go zero.RunAndBlock(&zero.Config{
		NickName:      []string{conf.BotConfig.Name},
		CommandPrefix: conf.BotConfig.CmdPrefix,
		SuperUsers:    conf.BotConfig.Admin,
		Driver: []zero.Driver{
			//driver.NewWebSocketServer(16, "ws://127.0.0.1:4322", ""),
			driver.NewWebSocketClient(conf.BotConfig.URL, conf.BotConfig.AccessToken),
		},
	}, nil)
	// 注册zerobot控制指令
	bot.DefineBotCommand()
	// 主线程锁死，并接收管道通信进行发送
	for {
		blrecMsg := <-network.LiveMsg
		log.Info("[channel] 收到用户ID (", blrecMsg[0], ") 的消息：", blrecMsg[1])
		uid, _ := strconv.ParseInt(blrecMsg[0], 10, 64)
		bot.SendMessageByUser(uid, blrecMsg[1])
	}
}
