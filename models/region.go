package models

import (
	db "github.com/Fengxq2014/sel/database"
)

type Provinces struct {
	Id         int `json:"id" form:"id"`
	Provinceid int `json:"provinceid" form:"provinceid"`
	Province   int `json:"province" form:"province"`
}

type Cities struct {
	Id         int `json:"id" form:"id"`
	Cityid     int `json:"cityid" form:"cityid"`
	City       int `json:"city" form:"city"`
	Provinceid int `json:"provinceid" form:"provinceid"`
}

// GetProvinces 获取省、直辖市信息
func GetProvinces() (provinces []Provinces, err error) {
	err = db.Engine.Find(&provinces)
	return provinces, err
}

// GetCities 获取地级市信息
func GetCities(provinceid int) (cities Cities, err error) {
	err = db.Engine.Where("provinceid=?", provinceid).Find(&cities)
	return cities, err
}
