package models

import (
	db "../database"
)

type Evaluation struct {
	Evaluation_id int    `json:"evaluation_id" form:"evaluation_id"`
	Name          string `json:"name" form:"name"`
	Category      string `json:"category" form:"category"`
	User_access   int    `json:"user_access" form:"user_access"`
	Abstract      string `json:"abstract" form:"abstract"`
	Details       string `json:"details" form:"details"`
	Price         int    `json:"price" form:"price"`
	Page_number   int    `json:"page_number" form:"page_number"`
	Person_count  int    `json:"person_count" form:"person_count"`
	Picture       string `json:"picture" form:"picture"`
	Sample_report string `json:"sample_report" form:"sample_report"`
}

// GetEvaluation 获取测评列表
func (e *Evaluation) GetEvaluation() (evaluation Evaluation, err error) {
	err = db.SqlDB.QueryRow("SELECT * FROM user where openid = ?", e.User_access).Scan(&evaluation.Evaluation_id)
	if err != nil {
		return evaluation, err
	}
	return evaluation, err
}
