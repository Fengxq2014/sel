package database

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
)

type config struct {
	Sqlname string
	Mysql   string
}

var conf *config
var SqlDB *sql.DB

func init() {
	var err error
	conf, err = readFile("app.conf")

	SqlDB, err = sql.Open(conf.Sqlname, conf.Mysql)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = SqlDB.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}
}

/*
获取配置文件对应字段值
*/
func readFile(filename string) (*config, error) {
	var conf = new(config)
	pwd, _ := os.Getwd()
	s := filepath.Join(pwd, "conf", filename)
	bytes, err := ioutil.ReadFile(s)
	if err != nil {
		log.Fatal("ReadFile: " + err.Error())
		return nil, err
	}
	if err := json.Unmarshal(bytes, &conf); err != nil {
		log.Fatal("Unmarshal: " + err.Error())
		return nil, err
	}

	return conf, nil
}
