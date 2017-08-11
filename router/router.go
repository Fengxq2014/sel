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
	router.GET("/", apis.IndexApi)
	//微信授权
	router.GET("/oauth", apis.Page1Handler)
	router.GET("/oauth1", apis.Page2Handler)
	router.Any("/weixin", apis.WeixinHandler)
	//登录
	router.POST("/login", apis.Login)
	//添加家长儿童关系
	router.GET("/adduser", apis.AddUcAPI)
	//获取测评列表
	router.GET("/getevalutionlist", apis.QryEvaluation)
	//获取题目
	router.GET("/getevalution", apis.QryQuestion)
	//上传答案
	router.GET("/updateevalution", apis.UpAnswer)
	//获取验证码
	router.GET("/sendcode", apis.SendSMS)
	return router
}
