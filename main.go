package main

import (
	"github.com/gin-gonic/gin"
	"go_sample/handler"
	"go_sample/logkit"
)

func init() {
	// zapLog.InitLog()
	logkit.Init(logkit.EnableCaller(true))
}

func main() {
	r := gin.Default()
	r.GET("/ping", handler.LogHandler)
	r.Run()
}