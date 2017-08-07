package router

import (
	. "../apis"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	// router := gin.Default()
	router := gin.New()
	router.Use(gin.Logger())
	// router.Use(logger())
	router.GET("/", IndexApi)

	router.POST("/person", AddPersonApi)

	router.GET("/persons", GetPersonsApi)

	router.GET("/person/:id", GetPersonApi)

	router.PUT("/person/:id", ModPersonApi)

	router.DELETE("/person/:id", DelPersonApi)

	router.Any("/weixin",WeixinHandler)
	return router
}
