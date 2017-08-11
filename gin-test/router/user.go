package router

import (
	"../apis"

	"github.com/gin-gonic/gin"
)

// InitRouter
func InitRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/", apis.IndexApi)
	//进入SEL，登录
	router.POST("/login/:openid/*telno/*name/*Unionid", apis.Login)
	router.Any("/weixin", apis.WeixinHandler)
	//微信授权
	router.GET("/oauth", apis.Page1Handler)
	router.GET("/oauth1", apis.Page2Handler)
	//添加家长儿童关系
	router.GET("/user/:user_id/*child_id/*relation", apis.InsertChild)
	//获取测评列表
	router.GET("/evalution/:user_access", apis.QryEvaluation)
	//获取题目类别
	router.GET("/evalution/:evaluation_id/*user_id/*child_id", apis.QryQuestion)
	//上传答案
	router.GET("/evalution/:evaluation_id/*user_id/*child_id/*current_question_id/*text_result/*report_result/*answer", apis.UpAnswer)
	return router
}
