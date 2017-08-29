package apis

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Fengxq2014/sel/models"
	"github.com/gin-gonic/gin"
)

type evaluationContain struct {
	Category    string              `json:"category"`
	Evaluations []models.Evaluation `json:"data"`
}

// QryEvaluation 获取测评列表
func QryEvaluation(c *gin.Context) {
	list := []evaluationContain{}
	caccess := c.Query("user_access")
	if caccess == "" {
		c.Error(errors.New("参数为空"))
		return
	}
	id, err := strconv.Atoi(caccess)
	if err != nil {
		c.Error(errors.New("参数不合法"))
		return
	}
	p := models.Evaluation{User_access: id}
	evaluation, err := p.GetEvaluation()
	if err != nil {
		c.Error(errors.New("查询有误"))
		return
	}
	if len(evaluation) > 0 {
		for _, value := range evaluation {
			index := checkExistCategory(&list, value.Category)
			if index > -1 {
				list[index].Evaluations = append(list[index].Evaluations, value)
			} else {
				eva := evaluationContain{Category: value.Category}
				eva.Evaluations = append(eva.Evaluations, value)
				list = append(list, eva)
			}
		}
	}
	c.JSON(http.StatusOK, models.Result{Data: &list})
}

// QryQuestion 获取题目
func QryQuestion(c *gin.Context) {
	type param struct {
		Eid   int `form:"evaluation_id" binding:"required"` //测评ID
		UID   int `form:"user_id" binding:"required"`       //用户ID
		CiD   int `form:"child_id" binding:"required"`      //儿童ID
		Index int `form:"index"`                            //题目号
	}
	var queryStr param
	if c.BindQuery(&queryStr) != nil {
		c.Error(errors.New("参数为空"))
		return
	}
	var err error
	res := models.Result{}
	user_evaluation := models.User_evaluation{Evaluation_id: queryStr.Eid, User_id: queryStr.UID, Child_id: queryStr.CiD}
	ue, err := user_evaluation.GetEvaluation()
	if ue.Current_question_id == -1 && err == nil {
		res.Res = 1
		res.Msg = "当前题目已经做完！"
		res.Data = nil
		c.JSON(http.StatusOK, res)
		return
	}
	var question models.Question
	if queryStr.Index > 0 {
		question, err = models.GetQuestionByIndex(queryStr.Eid, queryStr.Index)
	} else {
		question, err = models.GetQuestion(queryStr.Eid, queryStr.UID, queryStr.CiD)
	}

	if err != nil {
		res.Res = 1
		res.Msg = "获取题目失败" + err.Error()
		res.Data = nil
		c.JSON(http.StatusOK, res)
		return
	}
	res.Res = 0
	res.Msg = ""
	res.Data = question
	c.JSON(http.StatusOK, res)
}

// UpAnswer 上传答案
func UpAnswer(c *gin.Context) {
	type param struct {
		Eid      int    `form:"evaluation_id" binding:"required"`       //测评ID
		UID      int    `form:"user_id" binding:"required"`             //用户ID
		Cid      int    `form:"child_id" binding:"required"`            //儿童ID
		Cqid     int    `form:"current_question_id" binding:"required"` //当前测评题目，0为测评完成
		Tr       string `form:"text_result"`                            //测评描述
		Rr       string `form:"report_result"`                          //测评报告路径
		Answer   string `form:"answer" binding:"required"`              //答案
		MaxIndex int    `form:"maxIndex" binding:"required"`            //题目总数
	}
	//测评ID
	var queryStr param
	if c.BindQuery(&queryStr) != nil {
		c.Error(errors.New("参数为空"))
		return
	}
	if queryStr.Cqid == queryStr.MaxIndex {
		queryStr.Cqid = -1
		models.UpPersonCount()
	}
	err := models.UpdateUserAnswer(queryStr.Eid, queryStr.UID, queryStr.Cid, queryStr.Cqid, queryStr.Tr, queryStr.Rr, queryStr.Answer)
	res := models.Result{}
	if err != nil {
		res.Res = 1
		res.Msg = err.Error()
		res.Data = nil
		c.JSON(http.StatusOK, res)
		return
	}
	res.Res = 0
	res.Msg = ""
	res.Data = ""
	c.JSON(http.StatusOK, res)
}

func checkExistCategory(list *[]evaluationContain, category string) int {
	for index, value := range *list {
		if category == value.Category {
			return index
		}
	}
	return -1
}

// QryMyEvaluation 查询本人测评
func QryMyEvaluation(c *gin.Context) {
	type param struct {
		UID int `form:"user_id" binding:"required"` //用户ID
	}
	//测评ID
	var queryStr param
	if c.BindQuery(&queryStr) != nil {
		c.Error(errors.New("参数为空"))
		return
	}

	evaluation, err := models.GetMyEvaluation(queryStr.UID)
	if err != nil {
		c.Error(errors.New("查询有误"))
		return
	}
	c.JSON(http.StatusOK, models.Result{Data: evaluation})
}
