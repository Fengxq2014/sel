package models

import (
	db "github.com/Fengxq2014/sel/database"
)

type Course struct {
	Course_id    int    `json:"course_id" form:"course_id"`
	Name         string `json:"name" form:"name"`
	Category     string `json:"category" form:"category"`
	User_access  int    `json:"user_access" form:"user_access"`
	Valid_period int    `json:"valid_period" form:"valid_period"`
	Abstract     string `json:"abstract" form:"abstract"`
	Details      string `json:"details" form:"details"`
	Price        int    `json:"price" form:"price"`
	Person_count int    `json:"person_count" form:"person_count"`
	Picture      string `json:"picture" form:"picture"`
	Media        string `json:"media" form:"media"`
}

// GetCourse 获取课程列表
func (e *Course) GetCourse() (courses []Course, err error) {
	rows, err := db.SqlDB.Query("SELECT * FROM course where user_access = ?", e.User_access)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var course Course
		err = rows.Scan(&course.Course_id, &course.Name, &course.Category, &course.User_access, &course.Valid_period, &course.Abstract, &course.Details, &course.Price, &course.Person_count, &course.Picture, &course.Media)
		if err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}
	return courses, err
}

type User_course struct {
	User_course_id int    `json:"user_course_id" form:"user_course_id"`
	Course_id      string `json:"course_id" form:"course_id"`
	User_id        string `json:"user_id" form:"user_id"`
	Course_time    int    `json:"course_time" form:"course_time"`
}

// AddUsercourse 更新用户课程表
func (uc *User_course) AddUsercourse() (id int64, err error) {
	rs, err := db.SqlDB.Exec("INSERT INTO user_course(course_id,user_id,course_time) VALUES (?, ?, ?)", uc.Course_id, uc.User_id, uc.Course_time)
	if err != nil {
		return 0, err
	}
	id, err = rs.LastInsertId()
	return
}

// GetMyCourse 获取课程列表
func GetMyCourse(user_id int) (courses []Course, err error) {
	rows, err := db.SqlDB.Query("SELECT * FROM course where course_id in (select course_id from user_course where user_id=?)", user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var course Course
		err = rows.Scan(&course.Course_id, &course.Name, &course.Category, &course.User_access, &course.Valid_period, &course.Abstract, &course.Details, &course.Price, &course.Person_count, &course.Picture, &course.Media)
		if err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}
	return courses, err
}
