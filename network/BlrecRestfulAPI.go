package network

import (
	"BliveBrocastPostBot/data"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

// blrec restful api 操作
// 使用的api：
// GET /api/v1/tasks/{room_id}/cut 检查房间是否存在
// POST /api/v1/tasks/{room_id} 添加要监视的直播间
// DELETE /api/v1/tasks/{room_id} 删除要监视的直播间
// POST /api/v1/tasks/{room_id}/recorder/disable 取消自动录制
// POST /api/v1/tasks/{room_id}/recorder/enable 开启自动录制（管理员）

type blrecData struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type userInfo struct {
	UserInfo struct {
		Name string `json:"name"`
	} `json:"user_info"`
	TaskStatus struct {
		RecorderEnabled bool `json:"recorder_enabled"`
	} `json:"task_status"`
}

func blrecApiRequest(reqType string, apiUrl string, apiKey string) *http.Response {
	req, err := http.NewRequest(reqType, apiUrl, nil)
	if err != nil {
		log.Warn("[blrec api] 无法连接到blrecAPI！")
	}
	req.Header.Add("x-api-key", apiKey)
	req.Header.Add("accept", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Warn(err.Error())
	}
	return res
}

// BlrecCheckRoomAvailable 检查Blrec是否存在该房间
func BlrecCheckRoomAvailable(RoomID int64) bool {
	conf := <-data.ConfData
	data.ConfData <- conf
	var checkData blrecData
	apiUrl := conf.BlrecConfig.ApiURL + "/api/v1/tasks/" + strconv.FormatInt(RoomID, 10) + "/cut"
	taskResp := blrecApiRequest("GET", apiUrl, conf.BlrecConfig.APIKey)
	// 处理完成后关闭https连接
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Warn(err.Error())
		}
	}(taskResp.Body)
	err := json.NewDecoder(taskResp.Body).Decode(&checkData)
	if err != nil {
		log.Warn("[web] 解析BlrecJson出现错误！")
		return false
	}
	if checkData.Code != 404 {
		return true
	}
	return false
}

// BlrecAddRoom 在blrec中添加房间
func BlrecAddRoom(RoomID int64) {
	conf := <-data.ConfData
	data.ConfData <- conf
	apiUrl := conf.BlrecConfig.ApiURL + "/api/v1/tasks/" + strconv.FormatInt(RoomID, 10)
	taskResp := blrecApiRequest("POST", apiUrl, conf.BlrecConfig.APIKey)
	// 处理完成后关闭https连接
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Warn(err.Error())
		}
	}(taskResp.Body)
	all, err := io.ReadAll(taskResp.Body)
	if err != nil {
		return
	}
	log.Debug(string(all))
}

// BlrecDelRoom 在blrec删除房间
func BlrecDelRoom(RoomID int64) {
	conf := <-data.ConfData
	data.ConfData <- conf
	apiUrl := conf.BlrecConfig.ApiURL + "/api/v1/tasks/" + strconv.FormatInt(RoomID, 10)
	taskResp := blrecApiRequest("DELETE", apiUrl, conf.BlrecConfig.APIKey)
	// 处理完成后关闭https连接
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Warn(err.Error())
		}
	}(taskResp.Body)
	all, err := io.ReadAll(taskResp.Body)
	if err != nil {
		return
	}
	log.Debug(string(all))
}

// BlrecStopAutoRecord 停止blrec自动录制
func BlrecStopAutoRecord(RoomID int64) {
	conf := <-data.ConfData
	data.ConfData <- conf
	apiUrl := conf.BlrecConfig.ApiURL + "/api/v1/tasks/" + strconv.FormatInt(RoomID, 10) + "/recorder/disable"
	taskResp := blrecApiRequest("POST", apiUrl, conf.BlrecConfig.APIKey)
	// 处理完成后关闭https连接
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Warn(err.Error())
		}
	}(taskResp.Body)
	all, err := io.ReadAll(taskResp.Body)
	if err != nil {
		return
	}
	log.Debug(string(all))
}

// BlrecStartAutoRecord 开启blrec自动录制
func BlrecStartAutoRecord(RoomID int64) {
	conf := <-data.ConfData
	data.ConfData <- conf
	apiUrl := conf.BlrecConfig.ApiURL + "/api/v1/tasks/" + strconv.FormatInt(RoomID, 10) + "/recorder/enable"
	taskResp := blrecApiRequest("POST", apiUrl, conf.BlrecConfig.APIKey)
	// 处理完成后关闭https连接
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Warn(err.Error())
		}
	}(taskResp.Body)
	all, err := io.ReadAll(taskResp.Body)
	if err != nil {
		return
	}
	log.Debug(string(all))
}
