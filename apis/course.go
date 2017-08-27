package apis

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Fengxq2014/aliyun/vod"
	"github.com/Fengxq2014/sel/conf"
	"github.com/Fengxq2014/sel/models"
	"github.com/gin-gonic/gin"
)

type courseContain struct {
	Category string          `json:"category"`
	Course   []models.Course `json:"data"`
}

// QryCourse 获取课程列表
func QryCourse(c *gin.Context) {
	list := []courseContain{}
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

	if len(course) > 0 {
		for _, value := range course {
			index := checkExistCategorys(&list, value.Category)
			if index > -1 {
				list[index].Course = append(list[index].Course, value)
			} else {
				eva := courseContain{Category: value.Category}
				eva.Course = append(eva.Course, value)
				list = append(list, eva)
			}
		}
	}

	c.JSON(http.StatusOK, models.Result{Data: &list})
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

func checkExistCategorys(list *[]courseContain, category string) int {
	for index, value := range *list {
		if category == value.Category {
			return index
		}
	}
	return -1
}

// GetVideo 获取视频播放地址
func GetVideo(c *gin.Context) {
	res := models.Result{}
	Media := c.Query("media")
	playAuth, err := vod.NewAliyunVod(conf.Config.Access_key_id, conf.Config.Access_secret).GetVideoPlayAuth(Media)
	if err != nil {
		res.Res = 1
		res.Msg = err.Error()
		res.Data = nil
		c.JSON(http.StatusOK, res)
		return
	}
	res.Res = 0
	res.Msg = ""
	res.Data = playAuth.PlayAuth

	c.JSON(http.StatusOK, res)
}

// QryMyCourse 获取本人课程列表
func QryMyCourse(c *gin.Context) {
	list := []courseContain{}
	uid := c.Query("user_id")
	if uid == "" {
		c.Error(errors.New("参数为空"))
		return
	}
	id, err := strconv.Atoi(uid)
	if err != nil {
		c.Error(errors.New("参数不合法"))
		return
	}
	course, err := models.GetMyCourse(id)
	res := models.Result{}
	if err != nil {
		res.Res = 1
		res.Msg = "获取课程失败" + err.Error()
		res.Data = nil
		c.JSON(http.StatusOK, res)
		return
	}

	if len(course) > 0 {
		for _, value := range course {
			index := checkExistCategorys(&list, value.Category)
			if index > -1 {
				list[index].Course = append(list[index].Course, value)
			} else {
				eva := courseContain{Category: value.Category}
				eva.Course = append(eva.Course, value)
				list = append(list, eva)
			}
		}
	}

	c.JSON(http.StatusOK, models.Result{Data: &list})
}
