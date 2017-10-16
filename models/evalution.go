package models

import (
	"database/sql"
	"strconv"
	"time"

	db "github.com/Fengxq2014/sel/database"
)

type Evaluation struct {
	Evaluation_id       int       `json:"evaluation_id" form:"evaluation_id"`
	Name                string    `json:"name" form:"name"`
	Category            string    `json:"category" form:"category"`
	User_access         int       `json:"user_access" form:"user_access"`
	Abstract            string    `json:"abstract" form:"abstract"`
	Details             string    `json:"details" form:"details"`
	Price               float64   `json:"price" form:"price"`
	Page_number         int       `json:"page_number" form:"page_number"`
	Person_count        int       `json:"person_count" form:"person_count"`
	Picture             string    `json:"picture" form:"picture"`
	Sample_report       string    `json:"sample_report" form:"sample_report"`
	Key_name            string    `json:"key_name" form:"key_name"`
	MaxIndex            int64     `json:"maxIndex" form:"maxIndex" xorm:"-"`
	Current_question_id string    `json:"current_question_id" form:"current_question_id" xorm:"-"`
	Evaluation_time     time.Time `json:"evaluation_time" form:"evaluation_time" xorm:"-"`
	Child_id            int64     `json:"child_id" form:"child_id" xorm:"-"`
	User_evaluation_id  int64     `json:"user_evaluation_id" form:"user_evaluation_id" xorm:"-"`
}

// GetEvaluation 获取测评列表
func (e *Evaluation) GetEvaluation() (evaluations []Evaluation, err error) {
	err = db.Engine.Where("user_access=?", e.User_access).Find(&evaluations)
	var ev []Evaluation
	for _, values := range evaluations {
		question := Question{Evaluation_id: values.Evaluation_id}
		counts, _ := db.Engine.Count(&question)
		values.MaxIndex = counts
		ev = append(ev, values)
	}
	// rows, err := db.SqlDB.Query("SELECT * FROM evaluation where user_access = ?", e.User_access)
	// if err != nil {
	// 	return nil, err
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	var evaluation Evaluation
	// 	err = rows.Scan(&evaluation.Evaluation_id, &evaluation.Name, &evaluation.Category, &evaluation.User_access, &evaluation.Abstract, &evaluation.Details, &evaluation.Price, &evaluation.Page_number, &evaluation.Person_count, &evaluation.Picture, &evaluation.Sample_report)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	evaluations = append(evaluations, evaluation)
	// }
	return ev, err
}

type Question struct {
	Question_id    int    `json:"question_id" form:"question_id"`
	Evaluation_id  int    `json:"evaluation_id" form:"evaluation_id"`
	Question_index int    `json:"question_index" form:"question_index"`
	Content        string `json:"content" form:"content"`
	MaxIndex       int    `json:"maxIndex"  form:"maxIndex"`
	Answer         string `json:"answer"  form:"answer" xorm:"-"`
}

// GetQuestion 获取测评题目
// func GetQuestion(evaluation_id, user_id, child_id, currentIndex int) (question Question, err error) {
// 	var maxIndex int
// 	err = db.SqlDB.QueryRow("select Max(question_index) from question where evaluation_id=?", evaluation_id).Scan(&maxIndex)
// 	if err != nil {
// 		return question, err
// 	}
// 	question.MaxIndex = maxIndex
// 	if currentIndex > 0 {
// 		question.Question_index = currentIndex + 1
// 	} else {
// 		question.Question_index = 1
// 	}
// 	var Answer sql.NullString
// 	err = db.SqlDB.QueryRow("SELECT a.question_id,a.evaluation_id,a.question_index,a.content,b.answer FROM question a left join user_question b on b.question_id=a.question_index and b.user_id=? where evaluation_id = ? and question_index = ?", user_id, evaluation_id, question.Question_index).Scan(&question.Question_id, &question.Evaluation_id, &question.Question_index, &question.Content, &Answer)
// 	fmt.Println(user_id, evaluation_id, question.Question_index)
// 	if Answer.Valid {
// 		question.Answer = Answer.String
// 	}
// 	if err != nil {
// 		return question, err
// 	}
// 	return question, err
// }

// GetQuestionByIndex 根据index获取题目
func GetQuestionByIndex(evaluation_id, index, userID, User_evaluation_id int) (question Question, err error) {
	total, err := db.Engine.Where("evaluation_id=?", evaluation_id).Count(&question)

	// var maxIndex int
	// err = db.SqlDB.QueryRow("select Max(question_index) from question where evaluation_id=?", evaluation_id).Scan(&maxIndex)
	if err != nil {
		return question, err
	}
	question.MaxIndex = int64TOint(total)
	if index > question.MaxIndex {
		index = question.MaxIndex
	}
	var Answer sql.NullString

	err = db.SqlDB.QueryRow("select a.question_id,a.evaluation_id,a.question_index,a.content,b.answer from question a left join user_question b on b.question_id=a.question_id and b.user_id=? and b.user_evaluation_id=? where a.evaluation_id =? and a.question_index=? ", userID, User_evaluation_id, evaluation_id, index).Scan(&question.Question_id, &question.Evaluation_id, &question.Question_index, &question.Content, &Answer)
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
	Text_result         string    `json:"text_result" form:"text_result" xorm:"null"`
	Report_result       string    `json:"report_result" form:"report_result" xorm:"null"`
	TypeId              string    `json:"TypeId" form:"TypeId" xorm:"-"`
}

type User_question struct {
	User_question_id   int    `json:"user_question_id" form:"user_question_id"`
	User_evaluation_id int    `json:"user_evaluation_id" form:"user_evaluation_id"`
	Question_id        int    `json:"question_id" form:"question_id"`
	Answer             string `json:"answer" form:"answer" xorm:"null"`
	User_id            int    `json:"user_id" form:"user_id"`
}

// UpdateUserAnswer 上传答案
func UpdateUserAnswer(Evaluation_id, User_id, Child_id, Current_question_id, MaxIndex, Question_id int, Answer string) (err error) {
	ue := User_evaluation{Evaluation_id: Evaluation_id, User_id: User_id, Child_id: Child_id, Current_question_id: Current_question_id}
	uq := User_question{User_evaluation_id: Evaluation_id, Question_id: Question_id, Answer: Answer, User_id: User_id}
	//user_evaluation 有记录
	if selue, err := ue.GetEvaluation(); err == nil && selue.Current_question_id != 0 && selue.Current_question_id != -1 {
		ue.Current_question_id = maxInt(selue.Current_question_id, ue.Current_question_id) //防止修改答案时改变当前题目序号
		// if ue.Current_question_id == MaxIndex {
		// 	ue.Current_question_id = -1
		// 	UpPersonCount(Evaluation_id)
		// }
		ue.User_evaluation_id = selue.User_evaluation_id
		_, err := ue.UpdateEvaluation()
		if err != nil {
			return err
		}
		uq.User_evaluation_id = selue.User_evaluation_id
		uqq, err := uq.QryQuestion()
		if uqq.User_question_id != 0 {
			_, err = uq.UpQuestion()
			return err
		} else {
			_, err := uq.AddQuestion()
			if err != nil {
				return err
			}
		}
	} else {
		//user_evaluation 无记录
		ues := User_evaluation{Evaluation_id: Evaluation_id, User_id: User_id, Child_id: Child_id, Current_question_id: 1}
		id, err := ues.AddEvaluation()
		if id < 1 && err != nil {
			return err
		}
		selues, err := ue.GetEvaluation()
		uq.User_evaluation_id = selues.User_evaluation_id
		id, err = uq.AddQuestion()
		if id < 1 && err != nil {
			return err
		}
	}

	return err
}

// GetEvaluation 查询用户测评表
func (ue *User_evaluation) GetEvaluation() (uevaluation User_evaluation, err error) {
	_, err = db.Engine.Where("evaluation_id=? and user_id=? and child_id=? and current_question_id!=-1 ORDER BY evaluation_time DESC", ue.Evaluation_id, ue.User_id, ue.Child_id).Get(&uevaluation)
	// err = db.SqlDB.QueryRow("SELECT user_evaluation_id,current_question_id FROM user_evaluation where evaluation_id=? and user_id=? and child_id=?", ue.Evaluation_id, ue.User_id, ue.Child_id).Scan(&evaluation.User_evaluation_id, &evaluation.Current_question_id)
	// if err != nil {
	// 	return evaluation, err
	// }
	return uevaluation, err
}

func (ue *User_evaluation) UpdateEvaluation() (id int64, err error) {
	id, err = db.Engine.Cols("current_question_id").Update(ue, &User_evaluation{User_evaluation_id: ue.User_evaluation_id})
	// stmt, err := db.SqlDB.Prepare("update user_evaluation set current_question_id=?,text_result=?,report_result=?")
	// if err != nil {
	// 	return 0, err
	// }
	// defer stmt.Close()
	// rs, err := stmt.Exec(ue.Current_question_id, ue.Text_result, ue.Report_result)
	// if err != nil {
	// 	return 0, err
	// }
	// _, err = rs.RowsAffected()
	// if err != nil {
	// 	return 0, err
	// }
	return id, err
}

func (ue *User_evaluation) AddEvaluation() (id int64, err error) {
	ue.Evaluation_time = time.Now()
	id, err = db.Engine.Insert(ue)
	// rs, err := db.SqlDB.Exec("INSERT INTO user_evaluation(evaluation_id, user_id,child_id,evaluation_time,current_question_id,text_result,report_result) VALUES (?, ?, ?, ?, ?, ?, ?)", ue.Evaluation_id, ue.User_id, ue.Child_id, time.Now(), ue.Current_question_id, ue.Text_result, ue.Report_result)
	// if err != nil {
	// 	return 0, err
	// }
	// id, err = rs.LastInsertId()
	return
}

func (uq *User_question) AddQuestion() (id int64, err error) {
	id, err = db.Engine.Insert(uq)
	// rs, err := db.SqlDB.Exec("INSERT INTO user_question(user_evaluation_id,question_id,answer,user_id) VALUES (?, ?, ?,?)", uq.User_evaluation_id, uq.Question_id, uq.Answer, uq.User_id)
	// if err != nil {
	// 	return 0, err
	// }
	// id, err = rs.LastInsertId()
	return
}

func (uq *User_question) UpQuestion() (id int64, err error) {
	id, err = db.Engine.Cols("answer").Update(uq, &User_question{User_evaluation_id: uq.User_evaluation_id, Question_id: uq.Question_id, User_id: uq.User_id})
	// stmt, err := db.SqlDB.Prepare("update user_question set answer=? where user_evaluation_id=? and question_id=? and user_id=?")
	// defer stmt.Close()
	// _, err = stmt.Exec(uq.Answer, uq.User_evaluation_id, uq.Question_id, uq.User_id)
	return id, err
}

func (uq *User_question) QryQuestion() (uqs User_question, err error) {
	_, err = db.Engine.Where("user_evaluation_id=? and question_id=? and user_id=?", uq.User_evaluation_id, uq.Question_id, uq.User_id).Get(&uqs)
	// err = db.SqlDB.QueryRow("SELECT user_question_id FROM user_question where user_evaluation_id=? and question_id=?", uq.User_evaluation_id, uq.Question_id).Scan(&uqs.User_question_id)
	return uqs, err
}

func max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func maxInt(x, y int) int {
	if x > y {
		return x
	}
	return y
}

// GetMyEvaluation 获取本人测评
func GetMyEvaluation(user_id int) (evaluations []Evaluation, err error) {
	rows, err := db.SqlDB.Query("SELECT a.*,b.user_evaluation_id,b.current_question_id,b.evaluation_time,b.child_id FROM evaluation a left join user_evaluation b on b.evaluation_id=a.evaluation_id where b.user_id=?", user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var evaluation Evaluation
		err = rows.Scan(&evaluation.Evaluation_id, &evaluation.Name, &evaluation.Category, &evaluation.User_access, &evaluation.Abstract, &evaluation.Details, &evaluation.Price, &evaluation.Page_number, &evaluation.Person_count, &evaluation.Picture, &evaluation.Sample_report, &evaluation.Key_name, &evaluation.User_evaluation_id, &evaluation.Current_question_id, &evaluation.Evaluation_time, &evaluation.Child_id)
		if err != nil {
			return nil, err
		}
		evaluations = append(evaluations, evaluation)
	}
	return evaluations, err
}

// int64装换成int
func int64TOint(id64 int64) (id int) {
	//int64到string
	idstring := strconv.FormatInt(id64, 10)
	//string到int
	id, _ = strconv.Atoi(idstring)
	return
}

// UpPersonCount 更新已测评人数
func UpPersonCount(evaluation_id int) (err error) {
	stmt, _ := db.SqlDB.Prepare("update evaluation set person_count = person_count + 1 where evaluation_id =? ")
	defer stmt.Close()
	_, err = stmt.Exec(evaluation_id)
	return
}

// QryUserEvaluation 根据evaluation_id，user_id，child_id查询user_evaluation_id
func (ue *User_evaluation) QryUserEvaluation() (result User_evaluation, err error) {
	if ue.TypeId == "0" {
		_, err = db.Engine.Where("evaluation_id=? and user_id=? and child_id=? and current_question_id!=-1", ue.Evaluation_id, ue.User_id, ue.Child_id).Get(&result)
		return
	}
	_, err = db.Engine.Where("evaluation_id=? and user_id=? and child_id=? and current_question_id=-1", ue.Evaluation_id, ue.User_id, ue.Child_id).Get(&result)
	return
}

// QryUserEvaluation 根据evaluation_id，user_id，child_id查询user_evaluation_id
func (ue *User_evaluation) QryEvaluation() (result User_evaluation, err error) {
	_, err = db.Engine.Where("user_evaluation_id=?", ue.User_evaluation_id).Get(&result)
	return
}

// QryEvaluation 根据evaluation_id查询keyname
func QryEvaluation(evaluation_id int) (result Evaluation, err error) {
	_, err = db.Engine.Where("evaluation_id=?", evaluation_id).Get(&result)
	return
}
