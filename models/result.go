package models

// Result api返回model
type Result struct {
	Res  int
	Msg  string
	Data interface{}
}

// config 配置文件
type config struct {
	Sqlname string
	Mysql   string
}