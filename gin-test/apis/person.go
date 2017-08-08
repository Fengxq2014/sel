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

//登录判断，第一次登录插入客户信息
func GetUserApi(c *gin.Context) {
	cid := c.Param("id")

	p := User{Unionid: cid}
	user, err := p.Get()
	if err != nil {
		p := User{Unionid: cid, Role: 0}
		ra, err := p.Add()
		if err != nil {
			log.Fatalln(err)
		}
		msg := fmt.Sprintf("insert successful %d", ra)
		c.JSON(http.StatusOK, gin.H{
			"msg": msg,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

//注册插入用户信息表
func AddUserApi(c *gin.Context) {
	cid := c.Param("id")
	phone := c.Param("phone")
	name := c.Param("name")
	p := User{Unionid: cid, Phone_number: phone, Name: name}
	ra, err := p.Insert()
	if err != nil {
		log.Fatalln(err)
	}
	msg := fmt.Sprintf("insert successful %d", ra)
	c.JSON(http.StatusOK, gin.H{
		"msg": msg,
	})
}
