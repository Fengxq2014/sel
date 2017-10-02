package router

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/Fengxq2014/sel/apis"
	"github.com/Fengxq2014/sel/tool"
	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由
func InitRouter() *gin.Engine {
	pwd, _ := os.Getwd()
	s := filepath.Join(pwd, "log", "server.log")
	myfile, _ := os.OpenFile(s, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	gin.DefaultWriter = io.MultiWriter(myfile, os.Stdout)
	router := gin.Default()
	router.Use(handleErrors)
	router.Static("/front", "./front")
	router.StaticFile("/MP_verify_wKkoD2xPfCrtcZer.txt", "./front/MP_verify_wKkoD2xPfCrtcZer.txt")

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
	//通过openid查询用户信息
	router.GET("/qryuser", apis.QryUserAPI)
	//查询儿童信息
	router.GET("/qrychild", apis.QryUcAPI)
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
	//获取课程列表
	router.GET("/QryCourse", apis.QryCourse)
	//更新用户课程表
	router.GET("/UpUserCouse", apis.UpUserCouse)
	//获取视频播放地址
	router.GET("/GetVideoPlayAuth", apis.GetVideo)
	//上传儿童头像
	router.GET("/UploadChildImg", apis.DownloadMedia)
	//查询本人测评
	router.GET("/QryMyEvaluation", apis.QryMyEvaluation)
	//查询本人课程
	router.GET("/QryMyCourse", apis.QryMyCourse)
	//插入视频播放记录
	router.GET("/VideoPlaybackRecord", apis.QryMyVideo)
	//查看报告
	router.GET("/QryReport", apis.QryReport)
	//生成支付订单
	router.GET("/wxPayOrder", apis.WxPayOrder)
	//微信支付回调
	//router.GET("/wxPayCallBack", apis.WxPayCallBack)
	return router
}

func handleErrors(c *gin.Context) {
	c.Next()
	errorToPrint := c.Errors.Last()
	if errorToPrint != nil {
		c.JSON(200, gin.H{
			"res":  500,
			"msg":  errorToPrint.Error(),
			"data": nil,
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
