package apis

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Fengxq2014/sel/models"
	"github.com/gin-gonic/gin"
)

// QryEvaluation 获取测评列表
func QryEvaluation(c *gin.Context) {
	caccess := c.Param("user_access")
	id, err := strconv.Atoi(caccess)
	p := models.Evaluation{User_access: id}
	evaluation, err := p.GetEvaluation()
	res := models.Result{}
	if err != nil {
		c.Error(errors.New("没有找到相关信息"))
		return
	}
	res.Res = 0
	res.Msg = ""
	res.Data = evaluation
	c.JSON(http.StatusOK, res)
}

// QryQuestion 获取题目
func QryQuestion(c *gin.Context) {
	type param struct {
		Eid int `form:"evaluation_id" binding:"required"` //测评ID
		UID int `form:"user_id" binding:"required"`       //用户ID
		CiD int `form:"child_id" binding:"required"`      //儿童ID
	}
	var queryStr param
	if c.BindQuery(&queryStr) != nil {
		c.Error(errors.New("参数为空"))
		return
	}
	evaluation, err := models.GetQuestion(queryStr.Eid, queryStr.UID, queryStr.CiD)
	res := models.Result{}
	if err != nil {
		res.Res = 1
		res.Msg = "获取题目失败"
		res.Data = nil
		c.JSON(http.StatusOK, res)
	}
	res.Res = 0
	res.Msg = ""
	res.Data = evaluation
	c.JSON(http.StatusOK, res)
}

// UpAnswer 上传答案
func UpAnswer(c *gin.Context) {
	type param struct {
		Eid    int    `form:"evaluation_id" binding:"required"`       //测评ID
		UID    int    `form:"user_id" binding:"required"`             //用户ID
		Cid    int    `form:"child_id" binding:"required"`            //儿童ID
		Cqid   int    `form:"current_question_id" binding:"required"` //当前测评题目，0为测评完成
		Tr     string `form:"text_result"`                            //测评描述
		Rr     string `form:"report_result"`                          //测评报告路径
		Answer string `form:"answer" binding:"required"`              //答案
	}
	//测评ID
	var queryStr param
	if c.BindQuery(&queryStr) != nil {
		c.Error(errors.New("参数为空"))
		return
	}
	err := models.UpdateUserAnswer(queryStr.Eid, queryStr.UID, queryStr.Cid, queryStr.Cqid, queryStr.Tr, queryStr.Rr, queryStr.Answer)
	res := models.Result{}
	if err != nil {
		res.Res = 1
		res.Msg = "更新答案错误"
		res.Data = nil
		c.JSON(http.StatusOK, res)
	}
	res.Res = 0
	res.Msg = ""
	res.Data = ""
	c.JSON(http.StatusOK, res)
}
