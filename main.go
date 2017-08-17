package main

import (
	"github.com/Fengxq2014/sel/router"

	"github.com/Fengxq2014/sel/conf"
	db "github.com/Fengxq2014/sel/database"
)

func main() {
	defer db.SqlDB.Close()
	router := router.InitRouter()
	router.Run(conf.Config.Port)
}
