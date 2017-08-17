package router

import (
	"errors"

	"github.com/Fengxq2014/sel/tool"
	// "time"
	"io"
	"os"
	"path/filepath"

	"github.com/Fengxq2014/sel/apis"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	pwd, _ := os.Getwd()
	s := filepath.Join(pwd, "log", "server.log")
	myfile, _ := os.OpenFile(s, os.O_APPEND|os.O_CREATE|os.O_RDWR, 066)
	gin.DefaultWriter = io.MultiWriter(myfile, os.Stdout)
	router := gin.Default()
	router.Use(handleErrors)
	router.Static("/front", "./front")
	router.GET("/", apis.IndexApi)
	// authorized := router.Group("/")
	// authorized.Use(jwtAuth)
	// {
	// 	authorized.GET("login", apis.Login)
	// }
	//微信授权
	router.GET("/oauth", apis.Page1Handler)
	router.GET("/oauth1", apis.Page2Handler)
	router.Any("/weixin", apis.WeixinHandler)
	//登录
	router.POST("/login", apis.Login)
	//添加家长儿童关系
	router.GET("/addchild", apis.AddUcAPI)
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

func handleErrors(c *gin.Context) {
	c.Next()
	errorToPrint := c.Errors.Last()
	if errorToPrint != nil {
		c.JSON(200, gin.H{
			"Res":  500,
			"Msg":  errorToPrint.Error(),
			"Data": nil,
		})
	}
}

func jwtAuth(c *gin.Context) {
	jwt := c.GetHeader("token")
	if jwt != "" {
		if result := tool.JWTVal(jwt); result {
			c.Next()
		}
	}
	c.AbortWithError(200, errors.New("jwt error"))
	return
}
