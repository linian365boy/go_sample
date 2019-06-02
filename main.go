package main

import (
	"github.com/gin-gonic/gin"
	"go_sample/handler"
	"go_sample/logkit"
	"go_sample/middleware"
)

var logger *logWrapper

func init() {
	// zapLog.InitLog()
	logger, _ = logkit.Init(logkit.EnableCaller(true))
}

func main() {
	r := gin.Default()
	r.Use(gin.Recovery())
	r.GET("/ping", middleware.RequestLogMiddleware(logger), handler.LogHandler)
	r.Run()
}