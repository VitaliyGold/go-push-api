package main

import (
	"github.com/gin-gonic/gin"
)

const FILE_PATH = "./logs";

func main() {

	r := gin.Default()

	t := &LogHandler{}

	r.POST("/shop/:shopID/notification", t.AddLog)
	r.POST("/shop/slow/:shopID/notification", t.SlowLogs)
	r.POST("/shop/externalError/:shopID/notification", t.GetExternalError)
	r.GET("/shop/:shopID/logs", t.GetLogs)
	r.DELETE("/shop/:shopID", t.RemoveLogs)

	r.Run();
}