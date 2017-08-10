package router

import (
	"../apis"

	"github.com/gin-gonic/gin"
)

// InitRouter
func InitRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/", apis.IndexApi)
	//进入SEL，验证微信账号
	router.GET("/user/:openid", apis.QryUserAPI)
	//进入SEL，登录
	router.POST("/login/:openid/:telno/:name/:Unionid", apis.Login)
	router.GET("/oauth1", apis.Page2Handler)
	//获取测评列表
	//添加家长儿童关系
	router.GET("/oauth", apis.Page1Handler)
	router.Any("/weixin", apis.WeixinHandler)
	//微信授权
	return router
}
