package apis

import (
	"fmt"
	"log"
	"net/http"

	. "../models"

	"github.com/gin-gonic/gin"
)

func IndexApi(c *gin.Context) {
	c.String(http.StatusOK, "It works")
}

// QryUserAPI 查询用户信息
func QryUserAPI(c *gin.Context) {
	cid := c.Param("openid")
	p := User{Openid: cid}
	user, err := p.GetUserByOpenid()
	res := Result{}
	if err != nil {
		res.Res = 1
		res.Msg = "没有该用户信息请登录！"
		res.Data = nil
		c.JSON(http.StatusOK, res)
	}
	res.Res = 0
	res.Msg = ""
	res.Data = user
	c.JSON(http.StatusOK, res)
}

// Login 登录判断
func Login(c *gin.Context) {
	cid := c.Param("openid")
	ctel := c.Param("telno")
	cname := c.Param("name")
	cunionid := c.Param("Unionid")

	p := User{Phone_number: ctel}
	_, err := p.GetUserByPhone()
	res := Result{}
	// 家长登录插入客户信息
	if err != nil {
		p := User{Unionid: cunionid, Role: 0, Name: cname, Openid: cid}
		ra, err := p.Insert()
		if err != nil {
			log.Println(err)
		}
		msg := fmt.Sprintf("insert successful %d", ra)
		res.Res = 1
		res.Msg = msg
		res.Data = nil
		c.JSON(http.StatusOK, res)
	} else {
		// 老师登录插入微信标识
		p := User{Unionid: cunionid, Phone_number: ctel, Openid: cid}
		ra, err := p.Update()
		if err != nil {
			res.Res = 1
			res.Msg = err.Error()
			res.Data = nil
			c.JSON(http.StatusOK, res)
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
	// cuid := c.Param("user_id")
	// ccid := c.Param("child_id")
	// cre := c.Param("relation")
	openid := c.Param("openid")
	p := User{Openid: openid}
	user, err := p.GetUserByOpenid()
	res := Result{}
	if err != nil {
		res.Res = 1
		res.Msg = "没有该用户信息请登录！"
		res.Data = nil
		c.JSON(http.StatusOK, res)
	}
	res.Res = 0
	res.Msg = ""
	res.Data = user
	c.JSON(http.StatusOK, res)
}
