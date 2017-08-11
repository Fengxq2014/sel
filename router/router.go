package router

import (
	"io"
	"os"
	"path/filepath"

	"../apis"

	"github.com/gin-gonic/gin"
)

// InitRouter
func InitRouter() *gin.Engine {
	pwd, _ := os.Getwd()
	s := filepath.Join(pwd, "log", "server.log")
	myfile, _ := os.OpenFile(s, os.O_APPEND|os.O_CREATE|os.O_RDWR, 066)
	gin.DefaultWriter = io.MultiWriter(myfile, os.Stdout)
	router := gin.Default()
	router.Use(logger())
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
