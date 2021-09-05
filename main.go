package main

import (
	"github.com/gin-gonic/gin"
	"login4/common"
)

//CREATE DATABASE `testlogin` CHARACTER SET 'utf8' COLLATE 'utf8_general_ci';


func main() {

	db := common.GetDB()
	defer  db.Close();
	r := gin.Default()
	r = CollectRoute(r)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

