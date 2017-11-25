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
	"github.com/gin-gonic/gin/binding"
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
	type param struct {
		CID int `form:"course_id" binding:"required"` //关联课程ID
		Uid int `form:"user_id" binding:"required"`   //关联用户ID
	}
	var queryStr param
	if c.ShouldBindWith(&queryStr, binding.Query) != nil {
		c.Error(errors.New("参数为空"))
		return
	}

	p := models.User_course{Course_id: queryStr.CID, User_id: queryStr.Uid, Course_time: time.Now()}
	id, err := p.AddUsercourse()
	res := models.Result{}
	if id != -1 && err != nil {
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
	type param struct {
		Media   string `form:"media" binding:"required"`    //视频ID
		Formats string `form:"formmats" binding:"required"` //视频流格式，多个用逗号分隔，支持格式mp4,m3u8,mp3
	}
	var queryStr param
	if c.ShouldBindWith(&queryStr, binding.Query) != nil {
		c.Error(errors.New("参数为空"))
		return
	}
	res := models.Result{}

	playInfo, err := vod.NewAliyunVod(conf.Config.Access_key_id, conf.Config.Access_secret).GetPlayInfo(queryStr.Media, queryStr.Formats, "")
	if err != nil {
		res.Res = 1
		res.Msg = err.Error()
		res.Data = nil
		c.JSON(http.StatusOK, res)
		return
	}
	res.Res = 0
	res.Msg = ""
	res.Data = map[string]string{"playAuth": playInfo.PlayInfoList.PlayInfo[0].PlayURL, "coverurl": playInfo.VideoBase.CoverURL}

	c.JSON(http.StatusOK, res)
}

// QryMyCourse 获取本人课程列表
func QryMyCourse(c *gin.Context) {
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

	c.JSON(http.StatusOK, models.Result{Data: course})
}

// QryMyVideo 插入视频播放记录
func QryMyVideo(c *gin.Context) {
	type param struct {
		CID int `form:"course_id" binding:"required"` //关联课程ID
		Uid int `form:"user_id" binding:"required"`   //关联用户ID
	}
	var queryStr param
	if c.ShouldBindWith(&queryStr, binding.Query) != nil {
		c.Error(errors.New("参数为空"))
		return
	}
	res := models.Result{}
	p := models.User_course{Course_id: queryStr.CID, User_id: queryStr.Uid, Course_time: time.Now()}
	uc, err := p.QryVideo()
	if uc.Course_id != 0 && uc.User_id != 0 {
		res.Res = 0
		res.Msg = "已有记录！"
		res.Data = nil
		c.JSON(http.StatusOK, res)
		return
	} else {
		_, err = p.InsertVideo()
		if err != nil {
			res.Res = 1
			res.Msg = "插入课程表失败！" + err.Error()
			res.Data = nil
			c.JSON(http.StatusOK, res)
			return
		}
	}

	c.JSON(http.StatusOK, res)
}

// QryPayCourse 课程是否已经支付
func QryPayCourse(c *gin.Context) {
	type param struct {
		CID int `form:"course_id" binding:"required"` //关联课程ID
		Uid int `form:"user_id" binding:"required"`   //关联用户ID
	}
	var queryStr param
	if c.ShouldBindWith(&queryStr, binding.Query) != nil {
		c.Error(errors.New("参数为空"))
		return
	}

	res := models.Result{}
	p := models.User_course{Course_id: queryStr.CID, User_id: queryStr.Uid}
	uc, err := p.QryVideo()

	if uc.Course_id != 0 && uc.User_id != 0 && err == nil {
		res.Res = 0
		res.Msg = "已支付！"
		res.Data = 0
		c.JSON(http.StatusOK, res)
		return
	}

	res.Res = 0
	res.Msg = "未支付！"
	res.Data = 1
	c.JSON(http.StatusOK, res)
}

// UpPayCourse 视频支付完成
func UpPayCourse(c *gin.Context) {
	type param struct {
		CID int `form:"course_id" binding:"required"` //关联课程ID
		Uid int `form:"user_id" binding:"required"`   //关联用户ID
	}
	var queryStr param
	if c.ShouldBindWith(&queryStr, binding.Query) != nil {
		c.Error(errors.New("参数为空"))
		return
	}

	res := models.Result{}
	p := models.User_course{Course_id: queryStr.CID, User_id: queryStr.Uid}
	id, err := p.AddUsercourse()

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
}

// GetCourseByID 根据id获取课程信息
func GetCourseByID(c *gin.Context) {
	type param struct {
		CID int `form:"course_id" binding:"required"` //课程ID
	}
	var queryStr param
	if c.ShouldBindWith(&queryStr, binding.Query) != nil {
		c.Error(errors.New("参数为空"))
		return
	}
	res := models.Result{}
	qry := models.Course{Course_id: queryStr.CID}
	course, err := qry.GetCourseByID()
	if err != nil {
		c.Error(errors.New("没有该课程"))
		return
	}
	res.Res = 0
	res.Msg = ""
	res.Data = course
	c.JSON(http.StatusOK, res)
}

// GetResource 获取课程资源
func GetResource(c *gin.Context) {
	type param struct {
		CID int `form:"course_id" binding:"required"` //课程ID
	}
	var queryStr param
	if c.ShouldBindWith(&queryStr, binding.Query) != nil {
		c.Error(errors.New("参数为空"))
		return
	}
	res := models.Result{}
	resource, err := models.QryResource(queryStr.CID)

	if err != nil {
		c.Error(errors.New("获取课程资源失败"))
		return
	}

	res.Res = 0
	res.Msg = ""
	res.Data = resource
	c.JSON(http.StatusOK, res)
}

// QryUserCourse 查看用户单个课程
func QryUserCourse(c *gin.Context) {
	type param struct {
		CID int `form:"course_id" binding:"required"` //课程ID
		UID int `form:"user_id" binding:"required"`   //用户ID
	}
	var queryStr param
	if c.ShouldBindWith(&queryStr, binding.Query) != nil {
		c.Error(errors.New("参数为空"))
		return
	}
	res := models.Result{}
	course, err := models.QryUserCourse(queryStr.CID, queryStr.UID)

	if err != nil {
		c.Error(errors.New("获取课程资源失败"))
		return
	}

	res.Res = 0
	res.Msg = ""
	res.Data = course.User_course_id
	c.JSON(http.StatusOK, res)
}
