package network

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// 通过uid获取roomid
// 使用API：
// https://api.bilibili.com/x/space/wbi/acc/info 获取用户个人信息
// 请求参数：mid、token、platform、web_location、w_rid(计算得)、wts(unix时间戳)
// https://api.live.bilibili.com/room/v1/Room/get_info
// 请求参数：room_id

type LiveInfo struct {
	RoomId int64 `json:"roomid"`
}

type UserInfo struct {
	Code int `json:"code"`
	Data struct {
		LiveRoom *LiveInfo `json:"live_room"`
	} `json:"data"`
}

type RoomInfo struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Message string `json:"message"`
	Data    struct {
		Uid int64 `json:"uid"`
	} `json:"data"`
}

func BiliBiliApiGetRequest(apiUrl string) *http.Response {
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		log.Warn(err.Error())
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) "+
		"Gecko/20100101 Firefox/114.0")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;"+
		"q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Add("Accept-Language", "zh-CN")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Warn(err.Error())
	}
	return res
}

func GetRoomID(uid int64) int64 {
	infoApi := "https://api.bilibili.com/x/space/wbi/acc/info?"
	param := map[string]string{
		"mid":          strconv.FormatInt(uid, 10),
		"token":        "",
		"platform":     "web",
		"web_location": "1550101",
	}
	paramStr := EncodeWbi(param)
	infoRes := BiliBiliApiGetRequest(infoApi + paramStr) // B站API访问
	var uInfo UserInfo
	var roomID int64
	err := json.NewDecoder(infoRes.Body).Decode(&uInfo)
	if err != nil {
		log.Warn(err.Error())
	}
	if uInfo.Code != 0 {
		roomID = 0
	} else {
		if uInfo.Data.LiveRoom == nil {
			roomID = 0
		} else {
			roomID = uInfo.Data.LiveRoom.RoomId
		}
	}
	return roomID
}

func GetUserID(roomid int64) int64 {
	infoApi := "https://api.live.bilibili.com/room/v1/Room/get_info?"
	infoRes := BiliBiliApiGetRequest(infoApi + "room_id=" + strconv.FormatInt(roomid, 10))
	var rInfo RoomInfo
	err := json.NewDecoder(infoRes.Body).Decode(&rInfo)
	if err != nil {
		log.Warn(err.Error())
	}
	return rInfo.Data.Uid
}
