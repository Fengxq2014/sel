package models

import (
	"errors"
	"time"

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
	Price        string `json:"price" form:"price"`
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

// GetCourseByID 根据id获取课程信息
func (e *Course) GetCourseByID() (courses Course, err error) {
	has, err := db.Engine.Where("course_id = ?", e.Course_id).Get(&courses)
	if !has {
		return courses, errors.New("没有该课程")
	}
	return
}

type User_course struct {
	User_course_id int       `json:"user_course_id" form:"user_course_id"`
	Course_id      int       `json:"course_id" form:"course_id"`
	User_id        int       `json:"user_id" form:"user_id"`
	Course_time    time.Time `json:"course_time" form:"course_time"`
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
	rows, err := db.SqlDB.Query("SELECT * FROM course where course_id in (select course_id from user_course where user_id=? ORDER BY course_time DESC)", user_id)
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

// InsertVideo 插入视频播放记录
func (uc User_course) InsertVideo() (id int64, err error) {
	id, err = db.Engine.Insert(&uc)
	return
}

// QryVideo 查询用户课程表
func (uc *User_course) QryVideo() (user_course User_course, err error) {
	_, err = db.Engine.Where("course_id=? and user_id=?", uc.Course_id, uc.User_id).Get(&user_course)
	return user_course, err
}

type Resource struct {
	Resource_id int    `json:"resource_id" form:"resource_id"`
	Name        string `json:"name" form:"name"`
	Type        int    `json:"type" form:"resource_type"`
	Url         string `json:"url" form:"url"`
	Free        int    `json:"free" form:"free" xorm"-"`
	Index       int    `json:"index" form:"index" xorm"-"`
}

type Cresource struct {
	Cresource_id int `json:"cresource_id" form:"cresource_id"`
	Course_id    int `json:"course_id" form:"course_id"`
	Resource_id  int `json:"resource_id" form:"resource_id"`
	Index        int `json:"index" form:"index"`
	Free         int `json:"free" form:"free"`
}

// QryResource 获取课程资源
func QryResource(course_id int) (resource []Resource, err error) {
	err = db.Engine.Join("left", "cresource", "cresource.resource_id=resource.resource_id").Where("cresource.course_id = ?", course_id).Find(&resource)
	return resource, err
}

// QryUserCourse 查看用户单个课程
func QryUserCourse(course_id, user_id int) (course User_course, err error) {
	_, err = db.Engine.Where("course_id=? and user_id=?", course_id, user_id).Get(&course)
	return course, err
}
