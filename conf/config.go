package conf

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type config struct {
	Sqlname           string
	Mysql             string
	WXAppID           string
	WXAppSecret       string
	WXOriID           string
	WXToken           string
	Oauth2RedirectURI string
	Port              string
	SmsID             string
	Access_key_id     string
	Access_secret     string
	Sign_name         string
	Cookietime        int
	Host              string
	Mch_id            string
	Mch_name          string
	Key               string
	CallBack          string
	Template_id       string
}

var Config *config

func init() {
	var err error
	Config, err = readFile("app.conf")
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
