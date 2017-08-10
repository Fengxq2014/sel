package models

import (
	"time"

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
	err = db.SqlDB.QueryRow("SELECT * FROM user where user_access = ?", e.User_access).Scan(&evaluation.Evaluation_id, &evaluation.Name, &evaluation.Category, &evaluation.User_access, &evaluation.Abstract, &evaluation.Details, &evaluation.Price, &evaluation.Page_number, &evaluation.Person_count, &evaluation.Picture, &evaluation.Sample_report)
	if err != nil {
		return evaluation, err
	}
	return evaluation, err
}

type Question struct {
	Question_id    int    `json:"question_id" form:"question_id"`
	Evaluation_id  int    `json:"evaluation_id" form:"evaluation_id"`
	Question_index int    `json:"question_index" form:"question_index"`
	Content        string `json:"content" form:"content"`
}

// GetQuestion 获取测评题目
func (q *Question) GetQuestion() (question Question, err error) {
	err = db.SqlDB.QueryRow("SELECT * FROM question where evaluation_id = ?", q.Evaluation_id).Scan(&question.Question_id, &question.Evaluation_id, &question.Question_index, &question.Content)
	if err != nil {
		return question, err
	}
	return question, err
}

type User_evaluation struct {
	User_evaluation_id  int       `json:"user_evaluation_id" form:"user_evaluation_id"`
	Evaluation_id       int       `json:"evaluation_id" form:"evaluation_id"`
	User_id             int       `json:"user_id" form:"user_id"`
	Child_id            int       `json:"child_id" form:"child_id"`
	Evaluation_time     time.Time `json:"evaluation_time" form:"evaluation_time"`
	Current_question_id int       `json:"current_question_id" form:"current_question_id"`
	Text_result         string    `json:"text_result" form:"text_result"`
	Report_result       string    `json:"report_result" form:"report_result"`
}

type User_question struct {
	User_question_id   int    `json:"user_question_id" form:"user_question_id"`
	User_evaluation_id int    `json:"user_evaluation_id" form:"user_evaluation_id"`
	Question_id        int    `json:"question_id" form:"question_id"`
	Answer             string `json:"answer" form:"answer"`
}

func UpdateUserAnswer(Evaluation_id, User_id, Child_id, Current_question_id int, Text_result, Report_result, Answer string) (err error) {
	ue := User_evaluation{Evaluation_id: Evaluation_id, User_id: User_id, Child_id: Child_id}
	uq := User_question{User_evaluation_id: Evaluation_id, Question_id: Current_question_id, Answer: Answer}
	//user_evaluation 有记录
	if uee, err := ue.GetEvaluation(); err != nil && uee.Child_id != 0 {
		err := ue.UpdateEvaluation()
		if err != nil {
			return err
		}
		id, err := uq.AddQuestion()
		if id < 1 && err != nil {
			return err
		}
	} else {
		//user_evaluation 无记录
		ues := User_evaluation{Evaluation_id: Evaluation_id, User_id: User_id, Child_id: Child_id}
		id, err := ues.AddEvaluation()
		if id < 1 && err != nil {
			return err
		}
		id, err = uq.AddQuestion()
		if id < 1 && err != nil {
			return err
		}
	}

	return err
}

// GetEvaluation 获取用户测评表
func (ue *User_evaluation) GetEvaluation() (evaluation User_evaluation, err error) {
	err = db.SqlDB.QueryRow("SELECT User_evaluation_id FROM user_evaluation where evaluation_id=? and user_id=? and child_id=?", ue.Evaluation_id, ue.User_id, ue.Child_id).Scan(&evaluation.User_evaluation_id)
	if err != nil {
		return evaluation, err
	}
	return evaluation, err
}

func (ue *User_evaluation) UpdateEvaluation() (err error) {
	stmt, err := db.SqlDB.Prepare("update user_evaluation set current_question_id=?,text_result=?,report_result=?")
	defer stmt.Close()
	if err != nil {
		return err
	}
	rs, err := stmt.Exec(ue.Current_question_id, ue.Text_result, ue.Report_result)
	if err != nil {
		return err
	}
	_, err = rs.RowsAffected()
	if err != nil {
		return err
	}
	return nil
}

func (ue *User_evaluation) AddEvaluation() (id int64, err error) {
	rs, err := db.SqlDB.Exec("INSERT INTO user_evaluation(evaluation_id, user_id,child_id,evaluation_time,current_question_id,text_result,report_result) VALUES (?, ?, ?, ?, ?, ?, ?)", ue.Evaluation_id, ue.User_id, ue.Child_id, ue.Evaluation_time, ue.Current_question_id, ue.Text_result, ue.Report_result)
	if err != nil {
		return 0, err
	}
	id, err = rs.LastInsertId()
	return
}

func (uq *User_question) AddQuestion() (id int64, err error) {
	rs, err := db.SqlDB.Exec("INSERT INTO user_question(user_evaluation_id,question_id,answer) VALUES (?, ?, ?)", uq.User_evaluation_id, uq.Question_id, uq.Answer)
	if err != nil {
		return 0, err
	}
	id, err = rs.LastInsertId()
	return
}