package network

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

// APIList：
// /api/v1/system/healthcheck

func HealthCheck(writer http.ResponseWriter, _ *http.Request) {
	_, err := writer.Write([]byte(string("{\"code\":0,\"msg\":\"ok\"}")))
	if err != nil {
		log.Warn(err.Error())
		return
	}
}

func StartAPIListen(addr string) {
	log.Info("[RESTful API] Http服务开始监听地址：http://", addr)
	http.HandleFunc("/api/v1/system/healthcheck", HealthCheck)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("开启Http服务监听出现错误:", err.Error())
	}
}
