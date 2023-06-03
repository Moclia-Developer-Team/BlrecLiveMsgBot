package network

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

var LiveMsg = make(chan [2]string)

// Blrec默认JSON结构
type blrecJson struct {
	Id      string          `json:"id"`   // webhook ID
	Date    string          `json:"date"` // webhook发送的日期
	MsgType string          `json:"type"` // 消息类型
	Data    json.RawMessage `json:"data"` // 实际承载的消息数据
}

// 开播下播使用的Data结构
type liveBeginAndEnd struct {
	UserInfo struct {
		Name string `json:"name"` // 主播用户名
		Uid  int64  `json:"uid"`  // 主播的个人ID（不知道啥时候会达到需要大数类的长度）
	} `json:"user_info"`
	RoomInfo struct {
		Uid    int64  `json:"uid"`     // 主播个人ID
		RoomId int64  `json:"room_id"` // 主播房间号
		Cover  string `json:"cover"`   // 直播间封面
		Title  string `json:"title"`   // 直播间标题
	} `json:"room_info"`
}

// RecvBlrecWebhook 接收blrec的webhook信息并处理成变量，通过管道传出给机器人处理
func RecvBlrecWebhook(writer http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	var jsonData blrecJson
	err := json.NewDecoder(req.Body).Decode(&jsonData)
	if err != nil {
		log.Warn("解析BlrecJson出现错误！")
		return
	}
	switch jsonData.MsgType {
	case "LiveBeganEvent":
		var dat liveBeginAndEnd
		// 解析BlrecWebhook的json数据，并存入dat变量
		err := json.Unmarshal(jsonData.Data, &dat)
		// 错误处理
		if err != nil {
			log.Warn("[webhook] 解析BlrecJsonData时出现错误！位于LiveBeganEvent")
			return
		}
		showMsg := "您关注的主播【" + dat.UserInfo.Name + "】正在直播：\n" +
			dat.RoomInfo.Title +
			"\n[CQ:image,file=" + dat.RoomInfo.Cover + "]" +
			"\nhttps://live.bilibili.com/" + strconv.FormatInt(dat.RoomInfo.RoomId, 10)
		retuInfo := [2]string{strconv.FormatInt(dat.UserInfo.Uid, 10), showMsg} // 构建通过管道传输的内容
		LiveMsg <- retuInfo                                                     // 将用户id和消息内容通过管道传输到bot接口
		break
	case "LiveEndedEvent":
		var dat liveBeginAndEnd
		err := json.Unmarshal(jsonData.Data, &dat)
		if err != nil {
			log.Warn("[webhook] 解析BlrecJsonData时出现错误！位于LiveEndedEvent")
			return
		}
		showMsg := "您关注的主播【" + dat.UserInfo.Name + "】下播了"
		retuInfo := [2]string{strconv.FormatInt(dat.UserInfo.Uid, 10), showMsg} // 构建通过管道传输的内容
		LiveMsg <- retuInfo                                                     // 将用户id和消息内容通过管道传输到bot接口
		break
	default:
		log.Debug(string(jsonData.Data))
	}
	_, _ = io.WriteString(writer, "{\"code\":0,\"return\":\"Success\"}") // 返回给客户端的内容
}

// StartWeb 开启web服务
func StartWeb(addr string, postUrl string) {
	log.Info("[webhook] Http服务开始监听地址：http://", addr, postUrl)
	http.HandleFunc(postUrl, RecvBlrecWebhook)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("开启Http服务监听出现错误:", err.Error())
	}
}
