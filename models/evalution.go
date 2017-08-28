package models

import (
	"database/sql"
	"time"

	db "github.com/Fengxq2014/sel/database"
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
func (e *Evaluation) GetEvaluation() (evaluations []Evaluation, err error) {
	rows, err := db.SqlDB.Query("SELECT * FROM evaluation where user_access = ?", e.User_access)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var evaluation Evaluation
		err = rows.Scan(&evaluation.Evaluation_id, &evaluation.Name, &evaluation.Category, &evaluation.User_access, &evaluation.Abstract, &evaluation.Details, &evaluation.Price, &evaluation.Page_number, &evaluation.Person_count, &evaluation.Picture, &evaluation.Sample_report)
		if err != nil {
			return nil, err
		}
		evaluations = append(evaluations, evaluation)
	}
	return evaluations, err
}

type Question struct {
	Question_id    int    `json:"question_id" form:"question_id"`
	Evaluation_id  int    `json:"evaluation_id" form:"evaluation_id"`
	Question_index int    `json:"question_index" form:"question_index"`
	Content        string `json:"content" form:"content"`
	MaxIndex       int    `json:"maxIndex"  form:"maxIndex"`
	Answer         string `json:"answer"  form:"answer"`
}

// GetQuestion 获取测评题目
func GetQuestion(evaluation_id, user_id, child_id int) (question Question, err error) {
	ue := User_evaluation{Evaluation_id: evaluation_id, User_id: user_id, Child_id: child_id}
	var maxIndex int
	err = db.SqlDB.QueryRow("select Max(question_index) from question where evaluation_id=?", evaluation_id).Scan(&maxIndex)
	if err != nil {
		return question, err
	}
	question.MaxIndex = maxIndex
	if uee, err := ue.GetEvaluation(); err == nil {
		if uee.Current_question_id > 0 {
			question.Question_index = uee.Current_question_id
		} else {
			question.Question_index = 1
		}
	} else {
		question.Question_index = 1
	}
	err = db.SqlDB.QueryRow("SELECT question_id,evaluation_id,question_index,content FROM question where evaluation_id = ? and question_index = ?", evaluation_id, question.Question_index).Scan(&question.Question_id, &question.Evaluation_id, &question.Question_index, &question.Content)
	if err != nil {
		return question, err
	}
	return question, err
}

// GetQuestionByIndex 根据index获取题目
func GetQuestionByIndex(evaluation_id, index int) (question Question, err error) {
	var maxIndex int
	err = db.SqlDB.QueryRow("select Max(question_index) from question where evaluation_id=?", evaluation_id).Scan(&maxIndex)
	if err != nil {
		return question, err
	}
	question.MaxIndex = maxIndex
	if index > maxIndex {
		index = maxIndex
	}
	var Answer sql.NullString
	err = db.SqlDB.QueryRow("select a.evaluation_id,a.question_index,a.content,b.answer from question a left join user_question b on b.question_id=a.question_id where a.evaluation_id =? and a.question_index=?", evaluation_id, index).Scan(&question.Evaluation_id, &question.Question_index, &question.Content, &Answer)
	if Answer.Valid {
		question.Answer = Answer.String
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

// UpdateUserAnswer 上传答案
func UpdateUserAnswer(Evaluation_id, User_id, Child_id, Current_question_id int, Text_result, Report_result, Answer string) (err error) {
	ue := User_evaluation{Evaluation_id: Evaluation_id, User_id: User_id, Child_id: Child_id, Current_question_id: Current_question_id}
	uq := User_question{User_evaluation_id: Evaluation_id, Question_id: Current_question_id, Answer: Answer}
	//user_evaluation 有记录
	if _, err := ue.GetEvaluation(); err == nil {
		err := ue.UpdateEvaluation()
		if err != nil {
			return err
		}
		uqq, err := uq.QryQuestion()
		if uqq.User_question_id != 0 {
			err = uq.UpQuestion()
			return err
		} else {
			_, err := uq.AddQuestion()
			if err != nil {
				return err
			}
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
	err = db.SqlDB.QueryRow("SELECT user_evaluation_id,current_question_id FROM user_evaluation where evaluation_id=? and user_id=? and child_id=?", ue.Evaluation_id, ue.User_id, ue.Child_id).Scan(&evaluation.User_evaluation_id, &evaluation.Current_question_id)
	if err != nil {
		return evaluation, err
	}
	return evaluation, err
}

func (ue *User_evaluation) UpdateEvaluation() (err error) {
	stmt, err := db.SqlDB.Prepare("update user_evaluation set current_question_id=?,text_result=?,report_result=?")
	if err != nil {
		return err
	}
	defer stmt.Close()
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
	rs, err := db.SqlDB.Exec("INSERT INTO user_evaluation(evaluation_id, user_id,child_id,evaluation_time,current_question_id,text_result,report_result) VALUES (?, ?, ?, ?, ?, ?, ?)", ue.Evaluation_id, ue.User_id, ue.Child_id, time.Now(), ue.Current_question_id, ue.Text_result, ue.Report_result)
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

func (uq *User_question) UpQuestion() (err error) {
	stmt, err := db.SqlDB.Prepare("update user_question set answer=? where user_evaluation_id=? and question_id=?")
	defer stmt.Close()
	_, err = stmt.Exec(uq.Answer, uq.User_evaluation_id, uq.User_question_id)
	return err
}

func (uq *User_question) QryQuestion() (uqs User_question, err error) {
	err = db.SqlDB.QueryRow("SELECT user_question_id FROM user_question where user_evaluation_id=? and question_id=?", uq.User_evaluation_id, uq.Question_id).Scan(&uqs.User_question_id)
	return uqs, err
}

func max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

// GetMyEvaluation 获取本人测评
func GetMyEvaluation(user_id int) (evaluations []Evaluation, err error) {
	rows, err := db.SqlDB.Query("SELECT * FROM evaluation where evaluation_id in (select evaluation_id from user_evaluation where user_id=?)", user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var evaluation Evaluation
		err = rows.Scan(&evaluation.Evaluation_id, &evaluation.Name, &evaluation.Category, &evaluation.User_access, &evaluation.Abstract, &evaluation.Details, &evaluation.Price, &evaluation.Page_number, &evaluation.Person_count, &evaluation.Picture, &evaluation.Sample_report)
		if err != nil {
			return nil, err
		}
		evaluations = append(evaluations, evaluation)
	}
	return evaluations, err
}
