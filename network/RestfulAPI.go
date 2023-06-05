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
