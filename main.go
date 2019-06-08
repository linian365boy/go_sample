package main

import (
	"github.com/gin-gonic/gin"
	"go_sample/handler"
)

func main() {
	router := gin.New()

	router.Use(Loggers())

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, "Hello world!")
	})

	router.GET("/ping", myHandler)

	router.Run(":8080")
}

func myHandler(c *gin.Context){
	logger, err := handler.NewProductionLogger()
	if err != nil {
		panic(err)
	}
	logger.Info("heheda.......")
	c.JSON(200, "Fuck hehe.")
	logger.Info("heheda.......")
}

func Loggers() gin.HandlerFunc {
	logger, err := handler.NewProductionLogger()
	if err != nil {
		panic(err)
	}

	return handler.Logger(logger, handler.Options{
		RequestBodyLimit:  2000,
		RequestQueryLimit: 2000,
		ResponseBodyLimit: 2000,
	})
}