package network

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var MixinKey string

// getWbiKeys 访问B站API获取wbi_img下的img_key和sub_key
// API：https://api.bilibili.com/x/web-interface/nav
// return：img_key, sub_key
func getWbiKeys() (string, string) {
	// nav接口返回的json部分结构体，只保留需要的信息用以解析秘钥
	type NavData struct {
		Data struct {
			WbiImg struct {
				ImgUrl string `json:"img_url"`
				SubUrl string `json:"sub_url"`
			} `json:"wbi_img"`
		} `json:"data"`
	}
	// 通过http请求获取初始密钥
	navApi := "https://api.bilibili.com/x/web-interface/nav"
	navRes := BiliBiliApiGetRequest(navApi) // 访问b站API通用函数
	// 处理完成后关闭https连接
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Warn(err.Error())
		}
	}(navRes.Body)
	var navKey NavData
	// 将json解析到变量
	err := json.NewDecoder(navRes.Body).Decode(&navKey)
	if err != nil {
		log.Warn(err.Error())
	}
	// 将图片链接处理成秘钥
	subKey := strings.Split(
		strings.TrimPrefix(navKey.Data.WbiImg.SubUrl, "https://i0.hdslb.com/bfs/wbi/"), ".")[0]
	imgKey := strings.Split(
		strings.TrimPrefix(navKey.Data.WbiImg.ImgUrl, "https://i0.hdslb.com/bfs/wbi/"), ".")[0]
	return imgKey, subKey
}

// getMixinKey 将获取到的Key进行重排加密，得到加密数据用的秘钥
func getMixinKey(origin string) string {
	// 密钥置换表
	mixinKeyEncTab := []int{
		46, 47, 18, 2, 53, 8, 23, 32, 15, 50, 10, 31, 58, 3, 45, 35, 27, 43, 5, 49,
		33, 9, 42, 19, 29, 28, 14, 39, 12, 38, 41, 13, 37, 48, 7, 16, 24, 55, 40,
		61, 26, 17, 0, 1, 60, 51, 30, 4, 22, 25, 54, 21, 56, 59, 6, 63, 57, 62, 11,
		36, 20, 34, 44, 52,
	}
	// 对原始key根据置换表进行重排序
	var newKey string
	for _, table := range mixinKeyEncTab {
		newKey += string(origin[table])
	}
	return newKey[0:32]
}

// EncodeWbi 将传输数据和加密秘钥进行加密，得到加密完成的请求信息
func EncodeWbi(param map[string]string) string {
	// 获取时间戳并将时间戳添加到参数列表
	now := time.Now()
	param["wts"] = strconv.FormatInt(now.Unix(), 10)
	// 对参数列表以Key进行排序
	var keys []string
	for key := range param {
		keys = append(keys, key)
	}
	sort.Sort(sort.StringSlice(keys))
	// 将排序后的参数列表转换成urlQuery字符串，并对!'()*进行过滤
	regex := regexp.MustCompile("[!'\\(\\)*]")
	var urlQuery string
	for _, key := range keys {
		urlQuery += url.QueryEscape(key) + "=" + url.QueryEscape(regex.ReplaceAllString(param[key], "")) + "&"
	}
	urlQuery = strings.TrimSuffix(urlQuery, "&") // 清除尾随的&
	salt := urlQuery + MixinKey
	// 对加盐后的UrlQuery进行MD5运算
	saltByte := []byte(salt)
	WRidByte := md5.Sum(saltByte)
	WRid := fmt.Sprintf("%x", WRidByte)
	// 将w_rid加到query末尾，并返回
	return urlQuery + "&w_rid=" + WRid
}

func UpdateMixinKeyOnStartup() {
	log.Info("开始更新Bilibili主站验证密钥")
	imgKey, subKey := getWbiKeys()
	MixinKey = getMixinKey(imgKey + subKey)
	log.Info("[biliAuth] 更新MixinKey成功，新key：" + MixinKey)
}

func UpdateMixinKey() {
	for {
		log.Info("开始更新Bilibili主站验证密钥")
		now := time.Now()
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		timex := time.NewTimer(next.Sub(now))
		_ = <-timex.C
		imgKey, subKey := getWbiKeys()
		MixinKey = getMixinKey(imgKey + subKey)
		log.Info("[biliAuth] 更新MixinKey成功，新key：" + MixinKey)
	}
}
