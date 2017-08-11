package apis

import (
	"net/http"
	"strconv"

	. "../models"

	"github.com/gin-gonic/gin"
	log "../tool"
)

// QryUserAPI 查询用户信息
func QryEvaluation(c *gin.Context) {
	logger:=log.GetLogger()
	logger.Println()
	caccess := c.Param("user_access")
	id, err := strconv.Atoi(caccess)
	p := Evaluation{User_access: id}
	user, err := p.GetEvaluation()
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
