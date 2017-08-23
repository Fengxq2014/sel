package apis

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Fengxq2014/sel/models"
	"github.com/gin-gonic/gin"
)

// QryCourse 获取课程列表
func QryCourse(c *gin.Context) {
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
	p := models.Course{User_access: id}
	course, err := p.GetCourse()
	res := models.Result{}
	if err != nil {
		res.Res = 1
		res.Msg = "获取课程失败" + err.Error()
		res.Data = nil
		c.JSON(http.StatusOK, res)
		return
	}
	res.Res = 0
	res.Msg = ""
	res.Data = course

	c.JSON(http.StatusOK, res)
}

// UpUserCouse 更新用户课程表
func UpUserCouse(c *gin.Context) {
	Courseid := c.Query("course_id")
	Userid := c.Query("user_id")
	if Courseid == "" {
		c.Error(errors.New("参数为空"))
		return
	}
	if Userid == "" {
		c.Error(errors.New("参数为空"))
		return
	}
	p := models.User_course{Course_id: Courseid, User_id: Userid, Course_time: time.Now().Day()}
	id, err := p.AddUsercourse()
	res := models.Result{}
	if id != 1 && err != nil {
		res.Res = 1
		res.Msg = "更新用户课程表失败" + err.Error()
		res.Data = nil
		c.JSON(http.StatusOK, res)
		return
	}
	res.Res = 0
	res.Msg = ""
	res.Data = ""

	c.JSON(http.StatusOK, res)
}
