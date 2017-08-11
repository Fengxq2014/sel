package main

import (
	. "./router"

	db "./database"
)

func main() {
	defer db.SqlDB.Close()
	router := InitRouter()
	router.Run(":8080")
}
