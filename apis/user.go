package apis

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	. "../models"

	"github.com/gin-gonic/gin"
)

func IndexApi(c *gin.Context) {
	c.String(http.StatusOK, "It works")
}

// QryUserAPI 查询用户信息
func QryUserAPI(c *gin.Context) {
	cid := c.Query("openid")
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
	cid := c.Query("openid")
	ctel := c.Query("telno")
	cname := c.Query("name")
	cunionid := c.Query("Unionid")

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
	cuid := c.Query("user_id")
	uid, err := strconv.Atoi(cuid)
	ccid := c.Query("child_id")
	cid, err := strconv.Atoi(ccid)
	cre := c.Query("relation")
	re, err := strconv.Atoi(cre)
	gid := c.Query("gender")
	ggid, err := strconv.Atoi(gid)
	name := c.Query("name")
	bd := c.Query("birth_date")
	t, _ := time.Parse("2006-01-02", bd)
	err = InsertChild(uid, cid, re, ggid, name, t)
	res := Result{}
	if err != nil {
		res.Res = 1
		res.Msg = "没有该用户信息请登录！"
		res.Data = nil
		c.JSON(http.StatusOK, res)
	}
	res.Res = 0
	res.Msg = ""
	res.Data = ""
	c.JSON(http.StatusOK, res)
}
