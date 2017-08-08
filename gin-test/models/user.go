package models

import (
	db "../database"
)

type User struct {
	User_id      int64  `json:"user_id" form:"user_id"`
	Phone_number string `json:"phone_number" form:"phone_number"`
	Unionid      string `json:"unionid" form:"unionid"`
	Name         string `json:"name" form:"name"`
	Role         int    `json:"role" form:"role"`
}

// Get 通过微信微信身份标识获取客户信息
func (u *User) Get() (user User, err error) {
	err = db.SqlDB.QueryRow("SELECT * FROM user where unionid=?", u.Unionid).Scan(&user.User_id, &user.Phone_number, &user.Unionid, &user.Name, &user.Role)

	if err != nil {
		return user, err
	}

	return user, nil
}

// Add 第一次登录添加客户信息
func (u *User) Add() (id int64, err error) {
	rs, err := db.SqlDB.Exec("INSERT INTO user(user_id, unionid, role) VALUES (?, ?, ?)", u.User_id, u.Unionid, u.Role)

	if err != nil {
		return
	}

	id, err = rs.LastInsertId()
	return
}

// Insert 注册
func (u *User) Insert() (id int64, err error) {
	rs, err := db.SqlDB.Exec("INSERT INTO user(phone_number, name) VALUES (?, ?, ?)", u.Phone_number, u.Name)

	if err != nil {
		return
	}

	id, err = rs.LastInsertId()
	return
}
