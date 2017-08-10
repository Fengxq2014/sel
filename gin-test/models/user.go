package models

import (
	db "../database"
)

type User struct {
	User_id      int    `json:"user_id" form:"user_id"`
	Phone_number string `json:"phone_number" form:"phone_number"`
	Unionid      string `json:"unionid" form:"unionid"`
	Name         string `json:"name" form:"name"`
	Role         int    `json:"role" form:"role"`
	Openid       string `json:"openid" form:"openid"`
}

// GetUserByOpenid 通过微信微信身份标识获取客户信息
func (u *User) GetUserByOpenid() (user User, err error) {
	err = db.SqlDB.QueryRow("SELECT * FROM user where openid = ?", u.Openid).Scan(&user.User_id, &user.Phone_number, &user.Unionid, &user.Name, &user.Role, &user.Openid)

	if err != nil {
		return user, err
	}

	return user, nil
}

// GetUserByPhone 通过微信微信身份标识获取客户信息
func (u *User) GetUserByPhone() (user User, err error) {
	err = db.SqlDB.QueryRow("SELECT * FROM user where phone_number=?", u.Phone_number).Scan(&user.User_id, &user.Phone_number, &user.Unionid, &user.Name, &user.Role, &user.Openid)

	if err != nil {
		return user, err
	}

	return user, nil
}

// Insert 家长注册
func (u *User) Insert() (id int64, err error) {
	rs, err := db.SqlDB.Exec("INSERT INTO user(phone_number, unionid, name, role, openid) VALUES (?, ?, ?, ?, ?)", u.Phone_number, u.Unionid, u.Name, u.Role, u.Openid)

	if err != nil {
		return
	}

	id, err = rs.LastInsertId()
	return
}

// Update 老师更新
func (u *User) Update() (id int64, err error) {
	rs, err := db.SqlDB.Exec("update user set unionid=? openid=?", u.Unionid, u.Openid)

	if err != nil {
		return
	}

	id, err = rs.LastInsertId()
	return
}
