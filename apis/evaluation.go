package apis

import (
	"net/http"
	"strconv"

	. "../models"

	log "../tool"
	"github.com/gin-gonic/gin"
)

// QryEvaluation 获取测评列表
func QryEvaluation(c *gin.Context) {
	logger := log.GetLogger()
	logger.Println("ceshi")
	caccess := c.Param("user_access")
	id, err := strconv.Atoi(caccess)
	p := Evaluation{User_access: id}
	evaluation, err := p.GetEvaluation()
	res := Result{}
	if err != nil {
		res.Res = 1
		res.Msg = "没有该用户信息请登录！"
		res.Data = nil
		c.JSON(http.StatusOK, res)
	}
	res.Res = 0
	res.Msg = ""
	res.Data = evaluation
	c.JSON(http.StatusOK, res)
}

// QryQuestion 获取题目类别
func QryQuestion(c *gin.Context) {
	//测评ID
	eid := c.Query("evaluation_id")
	eeid, err := strconv.Atoi(eid)
	//用户ID
	uid := c.Query("user_id")
	uuid, err := strconv.Atoi(uid)
	//儿童ID
	cid := c.Query("child_id")
	ccid, err := strconv.Atoi(cid)

	evaluation, err := GetQuestion(eeid, uuid, ccid)
	res := Result{}
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
	//测评ID
	eid := c.Query("evaluation_id")
	eeid, err := strconv.Atoi(eid)
	//用户ID
	uid := c.Query("user_id")
	uuid, err := strconv.Atoi(uid)
	//儿童ID
	cid := c.Query("child_id")
	ccid, err := strconv.Atoi(cid)
	//儿童ID
	cqid := c.Query("current_question_id")
	ccqid, err := strconv.Atoi(cqid)
	//测评描述
	tid := c.Query("text_result")
	//测评报告路径
	rid := c.Query("report_result")
	//答案
	aid := c.Query("answer")
	err = UpdateUserAnswer(eeid, uuid, ccid, ccqid, tid, rid, aid)
	res := Result{}
	if err != nil {
		res.Res = 1
		res.Msg = "获取题目失败"
		res.Data = nil
		c.JSON(http.StatusOK, res)
	}
	res.Res = 0
	res.Msg = ""
	res.Data = ""
	c.JSON(http.StatusOK, res)
}
