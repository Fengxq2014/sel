package apis

import (
	"errors"
	"net/http"

	"github.com/Fengxq2014/sel/models"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// GetProvinces 获取省、直辖市信息
func GetProvinces(c *gin.Context) {
	provinces, err := models.GetProvinces()
	if err != nil {
		c.Error(errors.New("查询有误"))
		return
	}
	c.JSON(http.StatusOK, models.Result{Data: &provinces})
}

// GetCities 获取地级市信息
func GetCities(c *gin.Context) {
	type param struct {
		Provinceid int `form:"provinceid" binding:"required"` //省、直辖市ID
	}
	var queryStr param
	if c.ShouldBindWith(&queryStr, binding.Query) != nil {
		c.Error(errors.New("参数为空"))
		return
	}
	cities, err := models.GetCities(queryStr.Provinceid)
	if err != nil {
		c.Error(errors.New("查询有误"))
		return
	}
	c.JSON(http.StatusOK, models.Result{Data: &cities})
}
