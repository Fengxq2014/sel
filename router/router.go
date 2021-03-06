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
	//查询单个儿童信息
	router.GET("/qrysinglechild", apis.QrySingleChild)
	//查询家长儿童relation
	router.GET("/QryRelation", apis.QryRelation)
	//获取relation列表
	router.GET("/GetRelation", apis.GetRelation)
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
	//生成报告
	router.GET("/QryReport", apis.QryReport)
	//查看报告
	router.GET("/QryReports", apis.QryReports)
	//生成支付订单
	router.GET("/wxPayOrder", apis.WxPayOrder)
	//微信支付回调
	router.GET("/wxPayCallBack", apis.WxPayCallBack)
	//测评是否已经支付
	router.GET("/qrypayevalution", apis.QryPayEvalution)
	//课程是否已经支付
	router.GET("/qrypaycourse", apis.QryPayCourse)
	//测评支付完成
	router.GET("/uppayevalution", apis.UpPayEvalution)
	//视频支付完成
	router.GET("/uppaycourse", apis.UpPayCourse)
	//根据id获取课程信息
	router.GET("/GetCourseByID", apis.GetCourseByID)
	//根据id获取测评信息
	router.GET("/GetEvalutionByID", apis.QrySingleEvaluation)
	//获取课程资源
	router.GET("/GetResource", apis.GetResource)
	//获取省、直辖市信息
	router.GET("/GetProvinces", apis.GetProvinces)
	//获取地级市信息
	router.GET("/GetCities", apis.GetCities)
	//更新个人中心信息
	router.GET("/UpdateUser", apis.UpdateUser)
	//获取个人中心信息
	router.GET("/QryUser", apis.QryUser)
	//查看用户单个课程
	router.GET("/QryUserCourse", apis.QryUserCourse)
	//查询所属儿童已完成测评列表
	router.GET("/QryEvaluationByChildId", apis.QryEvaluationByChildId)
	//查询测评是否购买
	router.GET("/QryEvaluationGM", apis.QryEvaluationGM)
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
