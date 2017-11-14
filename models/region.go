package models

import (
	db "github.com/Fengxq2014/sel/database"
)

type Provinces struct {
	Id         string `json:"-" form:"id"`
	Provinceid string `json:"value" form:"provinceid"`
	Province   string `json:"name" form:"province"`
	Parent int `json:"parent" xorm:"-"`
}

type Cities struct {
	Id         string `json:"-" form:"id"`
	Cityid     string `json:"value" form:"cityid"`
	City       string `json:"name" form:"city"`
	Provinceid string `json:"parent" form:"provinceid"`
}

// GetProvinces 获取省、直辖市信息
func GetProvinces() (provinces []Provinces, err error) {
	err = db.Engine.Find(&provinces)
	return provinces, err
}

// GetCities 获取地级市信息
func GetCities(provinceid int) (cities []Cities, err error) {
	err = db.Engine.Where("provinceid=?", provinceid).Find(&cities)
	return cities, err
}
