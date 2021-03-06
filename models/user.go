package models

import (
	"time"

	db "github.com/Fengxq2014/sel/database"
)

type User struct {
	User_id       int    `json:"user_id" form:"user_id"`
	Phone_number  string `json:"phone_number" form:"phone_number"`
	Unionid       string `json:"unionid" form:"unionid"`
	Name          string `json:"name" form:"name"`
	Role          int    `json:"role" form:"role"`
	Openid        string `json:"openid" form:"openid"`
	Head_portrait string `json:"head_portrait" form:"head_portrait"`
	Nick_name     string `json:"nick_name" form:"nick_name"`
	Gender        string `json:"gender" form:"gender"`
	Birth_date    string `json:"birth_date" form:"birth_date"`
	Residence     string `json:"residence" form:"residence"`
	Child_id      string `json:"child_id" form:"child_id" xorm:"-"`
}

// GetUserByOpenid 通过微信身份标识获取客户信息
func (u *User) GetUserByOpenid() (user User, err error) {
	_, err = db.Engine.Join("left", "uc_relation", "uc_relation.user_id=user.user_id").Where("user.openid = ?", u.Openid).Get(&user)
	return user, err
}

// GetUser 获取客户信息
func (u *User) GetUser() (user User, err error) {
	_, err = db.Engine.Get(user)
	return user, err
}

// Insert 家长注册
func (u *User) Insert() (id int64, err error) {
	rs, err := db.SqlDB.Exec("INSERT INTO user(phone_number, unionid, nick_name, role, openid,head_portrait) VALUES (?, ?, ?, ?, ?,?)", u.Phone_number, u.Unionid, u.Name, 0, u.Openid, u.Head_portrait)

	if err != nil {
		return
	}

	id, err = rs.LastInsertId()
	return
}

// Update 老师更新
func (u *User) Update() (id int64, err error) {
	rs, err := db.SqlDB.Exec("update user set unionid=? ,openid=?", u.Unionid, u.Openid)
	if err != nil {
		return
	}
	id, err = rs.LastInsertId()
	return
}

type Child struct {
	Child_id      int64     `json:"child_id" form:"child_id"`
	Name          string    `json:"name" form:"name"`
	Gender        int       `json:"gender" form:"gender"`
	Birth_date    time.Time `json:"birth_date" form:"birth_date"`
	Head_portrait string    `json:"head_portrait" form:"head_portrait"`
	Relation      int       `json:"relation" form:"relation" xorm:"-"`
	User_id       int       `json:"user_id" form:"user_id" xorm:"-"`
}

type Uc_relation struct {
	Uc_relation_id int   `json:"uc_relation_id" form:"uc_relation_id"`
	User_id        int   `json:"user_id" form:"user_id"`
	Child_id       int64 `json:"child_id" form:"child_id"`
	Relation       int   `json:"relation" form:"relation"`
}

// InsertChild 插入儿童信息及儿童用户关联表
func InsertChild(user_id int, child_id int64, relation, Gender int, Head_portrait, Name string, Birth_date time.Time) (err error) {
	session := db.Engine.NewSession()
	defer session.Close()
	// add Begin() before any action
	err = session.Begin()
	child := Child{Child_id: child_id, Name: Name, Gender: Gender, Birth_date: Birth_date, Head_portrait: Head_portrait}
	_, err = session.Insert(&child)
	if err != nil {
		session.Rollback()
		return
	}
	uc_relation := Uc_relation{User_id: user_id, Child_id: child_id, Relation: relation}
	_, err = session.Insert(&uc_relation)
	if err != nil {
		session.Rollback()
		return
	}
	// add Commit() after all actions
	err = session.Commit()
	if err != nil {
		return
	}

	return err
}

// Getchild 查询儿童信息
func (uc *Uc_relation) Getchild() (child []Child, err error) {
	err = db.Engine.Join("left", "uc_relation", "uc_relation.child_id=child.child_id").Where("uc_relation.user_id=?", uc.User_id).Find(&child)
	return child, err
}

// Getchild 查询儿童信息
func GetChildById(child_id int) (child Child, err error) {
	_, err = db.Engine.Where("child_id=?", child_id).Get(&child)
	return child, err
}

// UpChild 更新儿童信息
func (child *Child) UpChild() (id int64, err error) {
	session := db.Engine.NewSession()
	defer session.Close()
	// add Begin() before any action
	err = session.Begin()

	_, err = session.Cols("name", "gender", "birth_date", "head_portrait").Update(child, &Child{Child_id: child.Child_id})
	if err != nil {
		session.Rollback()
		return 0, err
	}

	uc_relation := Uc_relation{Relation: child.Relation}
	_, err = session.Cols("relation").Update(uc_relation, &Uc_relation{Child_id: child.Child_id, User_id: child.User_id})
	if err != nil {
		session.Rollback()
		return 0, err
	}

	err = session.Commit()
	if err != nil {
		return 0, err
	}
	return
}

// UpdateUser 更新个人中心信息
func (u *User) UpdateUser() (id int64, err error) {
	rs, err := db.SqlDB.Exec("update user set name=?,gender=?,birth_date=?,residence=?,nick_name=? where user_id=?", u.Name, u.Gender, u.Birth_date, u.Residence, u.Nick_name, u.User_id)
	if err != nil {
		return
	}
	id, err = rs.LastInsertId()
	return
}

// QryUser 获取个人中心信息
func QryUser(User_id int) (user User, err error) {
	bl, err := db.Engine.Where("User_id=?", User_id).Get(&user)
	if !bl {
		return user, err
	}
	return user, err
}

// QryRelation 获取relation
func QryRelation(User_id, Child_id int) (uc Uc_relation, err error) {
	_, err = db.Engine.Where("user_id=? and child_id=?", User_id, Child_id).Get(&uc)
	return uc, err
}

// QrySingleChild 查询单个儿童信息
func QrySingleChild(Child_id, User_id int) (child Child, err error) {
	rows, err := db.SqlDB.Query("SELECT child.birth_date,child.child_id,child.gender,child.head_portrait,child.name,uc_relation.user_id,uc_relation.relation from child LEFT JOIN uc_relation on uc_relation.child_id=child.child_id	where child.child_id=? and uc_relation.user_id=?", Child_id, User_id)
	if err != nil {
		return child, err
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&child.Birth_date, &child.Child_id, &child.Gender, &child.Head_portrait, &child.Name, &child.User_id, &child.Relation)
		if err != nil {
			return child, err
		}
	}
	return child, err
}
