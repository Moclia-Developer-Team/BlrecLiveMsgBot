package data

import (
	"errors"
	"strconv"
	"strings"
)

// CheckMid 检查mid是否是纯数字
func CheckMid(mid string) bool {
	_, err := strconv.Atoi(mid)
	if err != nil {
		return false
	}
	return true
}

// CheckPrefix 匹配字符串的前缀列表
func CheckPrefix(str string, prefixList []string) (string, error) {
	for _, prefix := range prefixList {
		if strings.HasPrefix(str, prefix) {
			return prefix, nil
		}
	}
	return "", errors.New("无前缀被匹配")
}
