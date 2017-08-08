package router

import (
	. "../apis"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/", IndexApi)
	//进入SEL，验证微信账号
	router.GET("/persons", GetUserApi)

	router.Any("/weixin", WeixinHandler)

	router.GET("/page1", Page1Handler)
	router.GET("/page2", Page2Handler)

	return router
}
