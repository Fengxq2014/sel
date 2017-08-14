package models

import (
	"time"

	db "github.com/Fengxq2014/sel/database"
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

// GetUserByPhone 通过微信微信身份标识获取客户信息.
// user User struct
func (u *User) GetUserByPhone() (user User, err error) {
	err = db.SqlDB.QueryRow("SELECT * FROM user where phone_number=?", u.Phone_number).Scan(&user.User_id, &user.Phone_number, &user.Unionid, &user.Name, &user.Role, &user.Openid)

	if err != nil {
		return user, err
	}

	return user, nil
}

// Insert 家长注册
func (u *User) Insert() (id int64, err error) {
	rs, err := db.SqlDB.Exec("INSERT INTO user(phone_number, unionid, name, role, openid) VALUES (?, ?, ?, ?, ?)", u.Phone_number, u.Unionid, u.Name, 0, u.Openid)

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

type Child struct {
	Child_id   int `json:"child_id" form:"child_id"`
	Name       int `json:"name" form:"name"`
	Gender     int `json:"gender" form:"gender"`
	Birth_date int `json:"birth_date" form:"birth_date"`
}

type Uc_relation struct {
	Uc_relation_id int `json:"uc_relation_id" form:"uc_relation_id"`
	User_id        int `json:"user_id" form:"user_id"`
	Child_id       int `json:"child_id" form:"child_id"`
	Relation       int `json:"relation" form:"relation"`
}

// InsertChild 插入儿童信息及儿童用户关联表
func InsertChild(user_id, child_id, relation, Gender int, Name string, Birth_date time.Time) (err error) {
	tx, err := db.SqlDB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare("INSERT INTO child VALUES (?,?,?)")
	stmt1, err := tx.Prepare("INSERT INTO uc_relation(user_id, child_id, relation) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close() // danger!
	defer stmt1.Close()
	_, err = stmt.Exec(Name, Gender, Birth_date)
	_, err = stmt1.Exec(user_id, child_id, relation)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
