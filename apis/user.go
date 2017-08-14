package apis

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Fengxq2014/sel/tool"

	"github.com/Fengxq2014/aliyun_sms"
	"github.com/Fengxq2014/sel/models"
	"github.com/gin-gonic/gin"
	"github.com/goroom/rand"
)

func IndexApi(c *gin.Context) {
	type param struct {
		Uid string `form:"user_id"`
		// cid  int       `form:"child_id"`
		// re   int       `form:"relation"`
		// ggid int       `form:"gender"`
		T time.Time `form:"birth_date"`
	}
	eva := models.Evaluation{Name: "适应性行为测评",
		Category:      "SEL能力评估",
		User_access:   0,
		Abstract:      "这是简介",
		Details:       "这是详细说明",
		Price:         0,
		Page_number:   10,
		Person_count:  50,
		Picture:       "http://img4.imgtn.bdimg.com/it/u=2104185324,1359413794&fm=26&gp=0.jpg",
		Sample_report: "/root/evaluation_report",
	}
	id, err := eva.InsertEvaluation()
	tool.Debug(err)
	tool.Debug(id)
	var query param
	if c.BindQuery(&query) == nil {
		c.AbortWithError(200, errors.New("errorsss"))
	}
	res := models.Result{}
	c.JSON(200, res)
	c.String(http.StatusOK, "It works")
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
	cid := c.Query("openid")
	ctel := c.Query("telno")
	cname := c.Query("name")
	cunionid := c.Query("unionid")
	number := c.Query("number")
	if cid == "" || ctel == "" || cname == "" || number == "" {
		c.Error(errors.New("参数为空"))
		return
	}
	session, err := sessionStorage.Get(ctel)
	if err != nil {
		c.Error(err)
		return
	}
	if session.(string) != number {
		c.Error(errors.New("验证码错误！"))
		return
	}
	p := models.User{Phone_number: ctel}
	_, err = p.GetUserByPhone()
	// 家长登录插入客户信息
	if err != nil {
		p := models.User{Unionid: cunionid, Role: 0, Name: cname, Openid: cid}
		ra, err := p.Insert()
		if err != nil {
			c.Error(err)
			return
		}
		msg := fmt.Sprintf("insert successful %d", ra)
		res.Res = 1
		res.Msg = msg
		res.Data = nil
		c.JSON(http.StatusOK, res)
	} else {
		// 老师登录插入微信标识
		p := models.User{Unionid: cunionid, Phone_number: ctel, Openid: cid}
		ra, err := p.Update()
		if err != nil {
			c.Error(err)
			return
		}
		msg := fmt.Sprintf("insert successful %d", ra)
		res.Res = 0
		res.Msg = msg
		res.Data = nil
		c.JSON(http.StatusOK, res)
	}
}

// AddUcAPI 用户儿童关联
func AddUcAPI(c *gin.Context) {
	type param struct {
		UID  int    `form:"user_id" binding:"required"`
		Cid  int    `form:"child_id" binding:"required"`
		Re   int    `form:"relation" binding:"required"`
		Ggid int    `form:"gender" binding:"required"`
		Name string `form:"name" binding:"required"`
		T    string `form:"birth_date" binding:"required"`
	}
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
	err = models.InsertChild(queryStr.UID, queryStr.Cid, queryStr.Re, queryStr.Ggid, queryStr.Name, t)
	res := models.Result{}
	if err != nil {
		c.Error(errors.New("没有该用户信息请登录！"))
		return
	}
	res.Res = 0
	res.Msg = ""
	res.Data = ""
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
	aliyun_sms, err := aliyun_sms.NewAliyunSms("薄荷叶", "SMS_83955022", "LTAIfScyzpJdTAFi", "Kw5STaGOvayPhzGEr4nrsvzu4cKK0z")
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
	err = sessionStorage.Add(telno, no)
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
