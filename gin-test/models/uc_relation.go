package models

import (
	db "../database"
)

type Uc_relation struct {
	Uc_relation_id int `json:"uc_relation_id" form:"uc_relation_id"`
	User_id        int `json:"user_id" form:"user_id"`
	Child_id       int `json:"child_id" form:"child_id"`
	Relation       int `json:"relation" form:"relation"`
}

// Add 用户儿童关联表插入
func (u *Uc_relation) Add() (id int64, err error) {
	rs, err := db.SqlDB.Exec("INSERT INTO uc_relation(uc_relation_id, user_id, child_id,relation) VALUES (?, ?, ?,?)", u.Uc_relation_id, u.User_id, u.Child_id, u.Relation)

	if err != nil {
		return
	}

	id, err = rs.LastInsertId()
	return
}
