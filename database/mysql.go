package database

import (
	"database/sql"
	"log"

	"github.com/Fengxq2014/sel/conf"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var SqlDB *sql.DB
var Engine *xorm.Engine

func init() {
	var err error
	SqlDB, err = sql.Open(conf.Config.Sqlname, conf.Config.Mysql)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = SqlDB.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}
	Engine, err = xorm.NewEngine(conf.Config.Sqlname, conf.Config.Mysql)
}
