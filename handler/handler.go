package handler

import (
	"github.com/gin-gonic/gin"
	"go_sample/logkit"
)

func LogHandler(c *gin.Context) {
	logkit.Info("I am revice a message.")
	TestLog("Hello Boby.")
	c.JSON(200, gin.H{
		"messag": "pong",
	})
	logkit.Info("I am will response a message.")
}


