package apis

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin/binding"

	"github.com/Fengxq2014/sel/conf"

	"github.com/Fengxq2014/aliyun_sms"
	"github.com/Fengxq2014/sel/models"
	"github.com/Fengxq2014/sel/tool"
	"github.com/gin-gonic/gin"
	"github.com/goroom/rand"
)

func IndexApi(c *gin.Context) {
	runPrint("selreport", "52,fengtestdemo.pdf")
	c.String(http.StatusOK, "ok")
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
	if err != nil || user.Openid == "" {
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
		Cname         string `json:"name"`
		Cunionid      string `json:"unionid"`
		Number        string `json:"number" binding:"required"`
		Head_portrait string `json:"head_portrait"`
	}
	var postStr param
	if c.ShouldBindWith(&postStr, binding.JSON) != nil {
		c.Error(errors.New("参数为空"))
		return
	}
	sessionStorage.Delete(postStr.Ctel)
	p := models.User{Phone_number: postStr.Ctel}
	_, err := p.GetUser()
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
		UID           int    `form:"user_id" binding:"required"`
		Re            int    `form:"relation" binding:"required"`
		Ggid          int    `form:"gender" binding:"required"`
		Name          string `form:"name" binding:"required"`
		T             string `form:"birth_date" binding:"required"`
		CCID          int64  `form:"child_id"`
		Head_portrait string `form:"head_portrait"`
	}

	var queryStr param
	if c.ShouldBindWith(&queryStr, binding.Query) != nil {
		c.Error(errors.New("参数为空"))
		return
	}
	t, err := time.Parse("2006-01-02", queryStr.T)
	if err != nil {
		c.Error(errors.New("时间错误"))
		return
	}
	res := models.Result{}
	if queryStr.CCID != 0 {
		child := models.Child{Child_id: queryStr.CCID, Name: queryStr.Name, Gender: queryStr.Ggid, Birth_date: t, Head_portrait: queryStr.Head_portrait, Relation: queryStr.Re, User_id: queryStr.UID}
		_, err := child.UpChild()
		if err != nil {
			c.Error(errors.New("更新儿童信息失败！"))
			return
		}
		c.JSON(http.StatusOK, res)
		return
	}

	Cid := time.Now().Unix()
	err = models.InsertChild(queryStr.UID, Cid, queryStr.Re, queryStr.Ggid, queryStr.Head_portrait, queryStr.Name, t)
	if err != nil {
		c.Error(errors.New("插入儿童信息失败！"))
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

// QryUcAPI 查询儿童信息
func QryUcAPI(c *gin.Context) {
	cid := c.Query("user_id")
	if cid == "" {
		c.Error(errors.New("参数为空"))
		return
	}
	id, err := strconv.Atoi(cid)
	if err != nil {
		c.Error(errors.New("参数不合法"))
		return
	}
	p := models.Uc_relation{User_id: id}
	child, err := p.Getchild()
	res := models.Result{}
	if err != nil {
		res.Msg = "没有该儿童信息！"
		c.JSON(http.StatusOK, res)
		return
	}
	for _, value := range child {
		a := value.Birth_date.Format("2006-01-02")
		value.Birth_date, err = time.Parse("2006-01-02", a)
	}
	res.Res = 0
	res.Msg = ""
	res.Data = child
	c.JSON(http.StatusOK, res)
}

// UpdateUser 更新个人中心信息
func UpdateUser(c *gin.Context) {
	type param struct {
		Name       string `form:"name" binding:"required"`
		Gender     string `form:"gender" binding:"required"`
		Residence  string `form:"residence" binding:"required"`
		Birth_date string `form:"birth_date" binding:"required"`
		User_id    int    `form:"user_id" binding:"required"`
	}

	var queryStr param
	if c.ShouldBindWith(&queryStr, binding.Query) != nil {
		c.Error(errors.New("参数为空"))
		return
	}

	p := models.User{Name: queryStr.Name, Gender: queryStr.Gender, Residence: queryStr.Residence, Birth_date: queryStr.Birth_date, User_id: queryStr.User_id}

	id, err := p.UpdateUser()
	res := models.Result{}
	if id > 0 && err != nil {
		res.Msg = "更新个人中心信息失败！"
		c.JSON(http.StatusOK, res)
		return
	}
	res.Res = 0
	res.Msg = ""
	c.JSON(http.StatusOK, res)
}

// QryUser 获取个人中心信息
func QryUser(c *gin.Context) {
	type param struct {
		User_id int `form:"user_id" binding:"required"`
	}

	var queryStr param
	if c.ShouldBindWith(&queryStr, binding.Query) != nil {
		c.Error(errors.New("参数为空"))
		return
	}

	res := models.Result{}
	user, err := models.QryUser(queryStr.User_id)
	if err != nil {
		res.Msg = "更新个人中心信息失败！"
		c.JSON(http.StatusOK, res)
		return
	}
	res.Res = 0
	res.Msg = ""
	res.Data = user
	c.JSON(http.StatusOK, res)
}

// QryRelation 获取relation
func QryRelation(c *gin.Context) {
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
	user, err := models.QryRelation(queryStr.User_id, queryStr.Child_id)
	if err != nil {
		res.Msg = "查询家长儿童relation失败！"
		c.JSON(http.StatusOK, res)
		return
	}
	res.Res = 0
	res.Msg = ""
	res.Data = user
	c.JSON(http.StatusOK, res)
}

// 获取relation列表
func GetRelation(c *gin.Context) {
	res := models.Result{}
	res.Res = 0
	res.Msg = ""
	res.Data = map[string]string{"9": "未知", "1": "爸爸", "2": "妈妈", "3": "爷爷", "4": "奶奶", "5": "外公", "6": "外婆"}
	c.JSON(http.StatusOK, res)
	return
}

// QrySingleChild 查询单个儿童信息
func QrySingleChild(c *gin.Context) {
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
	child, err := models.QrySingleChild(queryStr.Child_id, queryStr.User_id)
	if err != nil {
		res.Msg = "查询单个儿童信息失败！"
		c.JSON(http.StatusOK, res)
		return
	}
	res.Res = 0
	res.Msg = ""
	res.Data = child
	c.JSON(http.StatusOK, res)
}
