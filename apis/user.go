package apis

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Fengxq2014/sel/conf"

	"github.com/Fengxq2014/aliyun_sms"
	"github.com/Fengxq2014/sel/models"
	"github.com/Fengxq2014/sel/tool"
	"github.com/gin-gonic/gin"
	"github.com/goroom/rand"
)

func IndexApi(c *gin.Context) {

}

// QryUserAPI 查询用户信息
func QryUserAPI(c *gin.Context) {
	cid := c.Query("openid")
	if cid == "" {
		c.Error(errors.New("参数为空"))
		return
	}
	p := models.User{Openid: cid}
	user, err := p.GetUserByOpenid()
	res := models.Result{}
	if err != nil {
		c.Error(errors.New("没有该用户信息请登录！"))
		return
	}
	res.Res = 0
	res.Msg = ""
	res.Data = user
	c.JSON(http.StatusOK, res)
}

// Login 登录判断.
func Login(c *gin.Context) {
	res := models.Result{}
	type param struct {
		ID            string `json:"openid" binding:"required"`
		Ctel          string `json:"telno" binding:"required"`
		Cname         string `json:"name" binding:"required"`
		Cunionid      string `json:"unionid"`
		Number        string `json:"number" binding:"required"`
		Head_portrait string `json:"head_portrait"`
	}
	var postStr param
	if c.BindJSON(&postStr) != nil {
		c.Error(errors.New("参数为空"))
		return
	}
	sessionStorage.Delete(postStr.Ctel)
	p := models.User{Phone_number: postStr.Ctel}
	_, err := p.GetUserByPhone()
	// 家长登录插入客户信息
	if err != nil {
		p := models.User{Unionid: postStr.Cunionid, Role: 0, Name: postStr.Cname, Openid: postStr.ID, Phone_number: postStr.Ctel}
		ra, err := p.Insert()
		if err != nil {
			c.Error(err)
			return
		}
		msg := fmt.Sprintf("insert successful %d", ra)
		res.Res = 0
		res.Msg = msg
		res.Data = nil
		c.JSON(http.StatusOK, res)
	} else {
		// 老师登录插入微信标识
		p := models.User{Unionid: postStr.Cunionid, Phone_number: postStr.Ctel, Openid: postStr.ID}
		ra, err := p.Update()
		if err != nil {
			c.Error(err)
			return
		}
		tool.Info("insert successful %d", ra)

		c.JSON(http.StatusOK, res)
	}
}

// AddUcAPI 用户儿童关联
func AddUcAPI(c *gin.Context) {
	type param struct {
		UID  int    `form:"user_id" binding:"required"`
		Re   int    `form:"relation" binding:"required"`
		Ggid int    `form:"gender" binding:"required"`
		Name string `form:"name" binding:"required"`
		T    string `form:"birth_date" binding:"required"`
	}
	Cid := time.Now().Unix()
	var queryStr param
	if c.BindQuery(&queryStr) != nil {
		c.Error(errors.New("参数为空"))
		return
	}
	t, err := time.Parse("2006-01-02", queryStr.T)
	if err != nil {
		c.Error(errors.New("时间错误"))
		return
	}
	err = models.InsertChild(queryStr.UID, Cid, queryStr.Re, queryStr.Ggid, queryStr.Name, t)
	res := models.Result{}
	if err != nil {
		c.Error(errors.New("没有该用户信息请登录！"))
		return
	}

	c.JSON(http.StatusOK, res)
}

// SendSMS 发送短信验证码
func SendSMS(c *gin.Context) {
	result := models.Result{Res: 1, Msg: "发送失败"}
	telno := c.Query("telno")
	if telno == "" {
		result.Msg = "参数为空"
		c.JSON(http.StatusOK, result)
		return
	}
	aliyun_sms, err := aliyun_sms.NewAliyunSms(conf.Config.Sign_name, conf.Config.SmsID, conf.Config.Access_key_id, conf.Config.Access_secret)
	if err != nil {
		tool.Error(err)
		result.Msg = err.Error()
		c.JSON(http.StatusOK, result)
		return
	}
	no := rand.String(4, rand.RST_NUMBER)
	tool.Debug("code:", no)
	err = aliyun_sms.Send(telno, `{"number":"`+no+`"}`)
	if err != nil {
		tool.Error(err)
		result.Msg = err.Error()
		c.JSON(http.StatusOK, result)
		return
	}
	err = sessionStorage.Set(telno, no)
	if err != nil {
		result.Msg = err.Error()
		c.JSON(http.StatusOK, result)
		return

	}
	result.Res = 0
	result.Msg = "成功"
	c.JSON(http.StatusOK, result)
	return
}
