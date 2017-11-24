package apis

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/gin-gonic/gin/binding"

	"github.com/Fengxq2014/sel/conf"
	"github.com/Fengxq2014/sel/models"
	"github.com/Fengxq2014/sel/tool"
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
	if c.ShouldBindWith(&queryStr, binding.Query) != nil {
		c.Error(errors.New("参数为空"))
		return
	}
	var err error
	res := models.Result{}
	user_evaluation := models.User_evaluation{Evaluation_id: queryStr.Eid, User_id: queryStr.UID, Child_id: queryStr.CiD}
	ue, err := user_evaluation.GetEvaluation()
	if err != nil {
		c.Error(err)
		return
	}
	var question models.Question
	if queryStr.Index > 0 {
		question, err = models.GetQuestionByIndex(queryStr.Eid, queryStr.Index, queryStr.UID, ue.User_evaluation_id)
	} else {
		if ue.Current_question_id > 0 {
			question, err = models.GetQuestionByIndex(queryStr.Eid, ue.Current_question_id+1, queryStr.UID, ue.User_evaluation_id)
		} else {
			question, err = models.GetQuestionByIndex(queryStr.Eid, 1, queryStr.UID, ue.User_evaluation_id)
		}
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
		Cqid     int    `form:"current_question_id" binding:"required"` //当前测评题目，-1为测评完成
		Tr       string `form:"text_result"`                            //测评描述
		Rr       string `form:"report_result"`                          //测评报告路径
		Answer   string `form:"answer" binding:"required"`              //答案
		MaxIndex int    `form:"maxIndex" binding:"required"`            //题目总数
		Qid      int    `form:"question_id" binding:"required"`         //测评题目ID
	}
	//测评ID
	var queryStr param
	if c.ShouldBindWith(&queryStr, binding.Query) != nil {
		c.Error(errors.New("参数为空"))
		return
	}

	err := models.UpdateUserAnswer(queryStr.Eid, queryStr.UID, queryStr.Cid, queryStr.Cqid, queryStr.MaxIndex, queryStr.Qid, queryStr.Answer)
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
	if c.ShouldBindWith(&queryStr, binding.Query) != nil {
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

// QryReport 生成报告
func QryReport(c *gin.Context) {
	type param struct {
		UEID   int    `form:"user_evaluation_id"`               //用户测评ID
		EID    int    `form:"evaluation_id" binding:"required"` //测评ID
		UID    int    `form:"user_id" binding:"required"`       //用户ID
		CID    int    `form:"child_id" binding:"required"`      //儿童ID
		TypeId string `form:"typeid"`                           //查看报告1；生成报告0
		OpenId string `form:"openid" binding:"required"`        //用户openid
	}
	var queryStr param
	if c.ShouldBindWith(&queryStr, binding.Query) != nil {
		c.Error(errors.New("参数为空"))
		return
	}
	res := models.Result{}
	ue := models.User_evaluation{Evaluation_id: queryStr.EID, User_id: queryStr.UID, Child_id: queryStr.CID, Current_question_id: -1, TypeId: queryStr.TypeId, User_evaluation_id: queryStr.UEID}

	use := models.User{Openid: queryStr.OpenId}
	uses, err := use.GetUserByOpenid()

	userEvaluation, err := ue.QryUserEvaluation()
	if err != nil {
		c.Error(err)
		return
	}
	evaluation, err := models.QryEvaluation(ue.Evaluation_id)
	if err != nil {
		c.Error(err)
		return
	}
	if len(userEvaluation.Report_result) <= 0 && queryStr.TypeId == "0" {
		idstring := strconv.Itoa(userEvaluation.User_evaluation_id)
		timetemp := time.Now().Format("20060102150405")
		fileName := evaluation.Key_name + "_" + idstring + "_" + timetemp + ".pdf"
		pdf := runPrint("selreport", idstring+","+fileName)
		if pdf {
			err := models.UpPersonCount(queryStr.EID)
			if err != nil {
				c.Error(err)
				return
			}
			_, err = ue.UpdateEvaluation()
			if err != nil {
				c.Error(err)
				return
			}

			err = TemplateMessage(queryStr.OpenId, conf.Config.Host+"/front/report/"+fileName, evaluation.Category, uses.Name)
			if err != nil {
				c.Error(err)
				return
			}

			ue.TypeId = "1"
			ue.User_evaluation_id = userEvaluation.User_evaluation_id
			userEvaluation, err = ue.QryUserEvaluation()
			if err != nil {
				c.Error(err)
				return
			}

			res.Res = 0
			res.Msg = ""
			res.Data = userEvaluation
			c.JSON(http.StatusOK, res)
			return
		}
		c.Error(errors.New("生成报告失败"))
		return
	}
	userEvaluation.Report_result = conf.Config.Host + userEvaluation.Report_result
	// err = TemplateMessage(queryStr.OpenId, "http://sel.bless-info.com"+userEvaluation.Report_result, evaluation.Category, uses.Name)
	// if err != nil {
	// 	c.Error(err)
	// 	return
	// }
	res.Res = 0
	res.Msg = ""
	res.Data = userEvaluation
	c.JSON(http.StatusOK, res)
	return
}

func runPrint(cmd string, args ...string) bool {
	os.Setenv("PATH", fmt.Sprintf("%s%c%s", "c:/sel/selreport", os.PathListSeparator, os.Getenv("PATH")))
	ecmd := exec.Command(cmd, args...)
	var errorout bytes.Buffer
	var out bytes.Buffer
	ecmd.Stdout = &out
	ecmd.Stderr = &errorout
	err := ecmd.Run()
	if err != nil {
		tool.Error(err)
	}
	if ecmd.ProcessState.Success() {
		return true
	}
	tool.Error(fmt.Sprintf("processstate:%v,out:%v,error:%v", ecmd.ProcessState, out.String(), errorout.String()))
	return false
}

// QryPayEvalution 测评是否已经支付
func QryPayEvalution(c *gin.Context) {
	type param struct {
		EID int `form:"evaluation_id" binding:"required"` //测评ID
		UID int `form:"user_id" binding:"required"`       //用户ID
		CID int `form:"child_id" binding:"required"`      //儿童ID
	}

	var queryStr param
	if c.ShouldBindWith(&queryStr, binding.Query) != nil {
		c.Error(errors.New("参数为空"))
		return
	}

	res := models.Result{}
	ue := models.User_evaluation{Evaluation_id: queryStr.EID, User_id: queryStr.UID, Child_id: queryStr.CID}
	ev, err := ue.GetEvaluation()

	if err == nil && ev.User_evaluation_id != 0 {
		res.Res = 0
		res.Msg = "已支付！"
		res.Data = 0
		c.JSON(http.StatusOK, res)
		return
	}

	res.Res = 0
	res.Msg = "未支付"
	res.Data = 1
	c.JSON(http.StatusOK, res)
	return
}

// UpPayEvalution 测评支付完成
func UpPayEvalution(c *gin.Context) {
	type param struct {
		EID int `form:"evaluation_id" binding:"required"` //测评ID
		UID int `form:"user_id" binding:"required"`       //用户ID
		CID int `form:"child_id" binding:"required"`      //儿童ID
	}

	var queryStr param
	if c.ShouldBindWith(&queryStr, binding.Query) != nil {
		c.Error(errors.New("参数为空"))
		return
	}

	res := models.Result{}
	ue := models.User_evaluation{Evaluation_id: queryStr.EID, User_id: queryStr.UID, Child_id: queryStr.CID}
	id, err := ue.AddEvaluation()

	if id != 1 && err != nil {
		res.Res = 1
		res.Msg = "更新用户课程表失败" + err.Error()
		res.Data = nil
		c.JSON(http.StatusOK, res)
		return
	}

	res.Res = 0
	res.Msg = ""
	res.Data = nil
	c.JSON(http.StatusOK, res)
	return
}

// QryReports 查看报告
func QryReports(c *gin.Context) {
	type param struct {
		UEID int `form:"user_evaluation_id" binding:"required"` //用户测评ID
		EID  int `form:"evaluation_id" binding:"required"`      //测评ID
	}
	var queryStr param
	if c.ShouldBindWith(&queryStr, binding.Query) != nil {
		c.Error(errors.New("参数为空"))
		return
	}
	res := models.Result{}
	ue := models.User_evaluation{Evaluation_id: queryStr.EID, User_evaluation_id: queryStr.UEID}

	userEvaluation, err := ue.QryEvaluation()
	if err != nil {
		c.Error(err)
		return
	}
	evaluation, err := models.QryEvaluation(ue.Evaluation_id)
	if err != nil {
		c.Error(err)
		return
	}
	res.Res = 0
	res.Msg = ""
	res.Data = map[string]string{"pdf": userEvaluation.Report_result, "details": evaluation.Details, "name": evaluation.Name, "reporttime": userEvaluation.Evaluation_time.Format("20060102150405"), "textResult": userEvaluation.Data_result}
	c.JSON(http.StatusOK, res)
	return
}

// QrySingleEvaluation 获取单个测评
func QrySingleEvaluation(c *gin.Context) {
	type param struct {
		EID int `form:"evaluation_id" binding:"required"` //测评ID
	}
	//测评ID
	var queryStr param
	if c.ShouldBindWith(&queryStr, binding.Query) != nil {
		c.Error(errors.New("参数为空"))
		return
	}
	res := models.Result{}
	evaluation, err := models.QryEvaluation(queryStr.EID)
	if err != nil {
		c.Error(err)
		return
	}

	res.Res = 0
	res.Msg = ""
	res.Data = evaluation
	c.JSON(http.StatusOK, res)
	return
}

// QryEvaluationByChildId 查询所属儿童测评列表
func QryEvaluationByChildId(c *gin.Context) {
	type param struct {
		User_id  int `form:"user_id" binding:"required"`
		Child_id int `form:"child_id" binding:"required"`
	}
	var queryStr param
	if c.ShouldBindWith(&queryStr, binding.Query) != nil {
		c.Error(errors.New("参数为空"))
		return
	}
	res := models.Result{}
	ue, err := models.QryEvaluationByChildId(queryStr.User_id, queryStr.Child_id)
	if err != nil {
		c.Error(err)
		return
	}

	res.Res = 0
	res.Msg = ""
	res.Data = ue
	c.JSON(http.StatusOK, res)
	return
}
