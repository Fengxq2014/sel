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
	Child_id      string `json:"child_id" form:"child_id"`
}

// GetUserByOpenid 通过微信微信身份标识获取客户信息
func (u *User) GetUserByOpenid() (user User, err error) {
	_, err = db.Engine.Join("left", "uc_relation", "uc_relation.user_id=user.user_id").Where("user.openid = ?", u.Openid).Get(&user)

	// var child_id sql.NullString
	// err = db.SqlDB.QueryRow("SELECT a.user_id,a.phone_number,a.name,a.role,a.head_portrait,b.child_id FROM user a left join uc_relation b on b.user_id=a.user_id where a.openid = ?", u.Openid).Scan(&user.User_id, &user.Phone_number, &user.Name, &user.Role, &user.Head_portrait, &child_id)
	// if child_id.Valid {
	// 	user.Child_id = child_id.String
	// }
	// if err != nil {
	// 	return user, err
	// }

	return user, err
}

// GetUser 获取客户信息
func (u *User) GetUser() (user User, err error) {
	_, err = db.Engine.Get(user)

	// err = db.SqlDB.QueryRow("SELECT * FROM user where phone_number=?", u.Phone_number).Scan(&user.User_id, &user.Phone_number, &user.Unionid, &user.Name, &user.Role, &user.Openid)
	// if err != nil {
	// 	return user, err
	// }

	return user, err
}

// Insert 家长注册
func (u *User) Insert() (id int64, err error) {
	rs, err := db.SqlDB.Exec("INSERT INTO user(phone_number, unionid, name, role, openid,head_portrait) VALUES (?, ?, ?, ?, ?,?)", u.Phone_number, u.Unionid, u.Name, 0, u.Openid, u.Head_portrait)

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
	// tx, err := db.SqlDB.Begin()
	// if err != nil {
	// 	return err
	// }
	// defer tx.Rollback()
	// stmt, err := tx.Prepare("INSERT INTO child VALUES (?,?,?,?,?)")
	// stmt1, err := tx.Prepare("INSERT INTO uc_relation(user_id, child_id, relation) VALUES (?, ?, ?)")
	// if err != nil {
	// 	return err
	// }
	// defer stmt.Close() // danger!
	// defer stmt1.Close()
	// _, err = stmt.Exec(child_id, Name, Gender, Birth_date, Head_portrait)
	// if err != nil {
	// 	return err
	// }
	// _, err = stmt1.Exec(user_id, child_id, relation)
	// if err != nil {
	// 	return err
	// }
	// err = tx.Commit()
	// if err != nil {
	// 	return err
	// }

	return err
}

// Getchild 查询儿童信息
func (uc *Uc_relation) Getchild() (child []Child, err error) {
	err = db.Engine.Join("left", "uc_relation", "uc_relation.child_id=child.child_id").Where("uc_relation.user_id=?", uc.User_id).Find(&child)

	// err = db.SqlDB.QueryRow("SELECT * FROM child where child_id in (select child_id from uc_relation where user_id=?)", uc.User_id).Scan(&child.Child_id, &child.Name, &child.Gender, &child.Birth_date, &child.Head_portrait)

	// if err != nil {
	// 	return child, err
	// }

	return child, err
}

// UpChild 更新儿童信息
func (child *Child) UpChild() (id int64, err error) {
	id, err = db.Engine.Cols("name", "gender", "birth_date", "head_portrait").Update(child, &Child{Child_id: child.Child_id})
	// rs, err := db.SqlDB.Exec("update child set name=? ,gender=? ,birth_date=? ,head_portrait=? where child_id=?", &child.Name, &child.Gender, &child.Birth_date, &child.Head_portrait, &child.Child_id)

	// if err != nil {
	// 	return
	// }

	// id, err = rs.LastInsertId()
	return
}
