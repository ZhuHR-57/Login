package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"math/rand"
	"net/http"
	"time"
)

//CREATE DATABASE `testlogin` CHARACTER SET 'utf8' COLLATE 'utf8_general_ci';

type User struct {
	gorm.Model
	Name string `gorm:"type:varchar(20); not null"`
	Telephone string `gorm:varchar(20);not null;unique`
	Password string `gorm:"size:255;not null"`
}

//打开连接池
func InitDB() *gorm.DB{
	driverName := "mysql"
	host := "localhost"
	port := "3306";
	database := "testlogin"
	username := "root"
	password := "root"
	charset := "utf8"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)

	db,err := gorm.Open(driverName,args)
	if err != nil {
		panic("failed to open database, err: "+err.Error())
	}

return db;
}

//随机字符
func RandomName(n int) string {
	var letters = []byte("diabvciwnvonwmxwomonrvbuebivnwoineo")
	result := make([]byte, n)

	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}

	return string(result)
}

func main() {
	db := InitDB()
	defer  db.Close();
	r := gin.Default()
	r.POST("/api/auth/register", func(ctx *gin.Context) {
		/*
			获取参数
		*/
		name := ctx.PostForm("name")
		telephone := ctx.PostForm("telephone")
		password := ctx.PostForm("password")

		/*
			数据验证
		*/
		if len(password) < 6 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码不能少于6位"})
			return
		}
		if len(telephone) != 11 {
			//gin.H = map[string]interface{}
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位", "number": len(telephone)})
			return
		}
		//名称没有传就会给一个10位的随机字符串
		if len(name) == 0 {
			name = RandomName(10)
		}

		//判断手机号是否存在
		if isTelephoneExist(db,telephone) {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户存在", "number": len(telephone)})
			return
		}

		//创建用户
		newUser := User{
			Name: name,
			Telephone: telephone,
			Password: password,
		}
		db.Create(&newUser);
		//返回结果
		ctx.JSON(200, gin.H{
			"message": "注册成功",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func isTelephoneExist(db *gorm.DB,telephone string) bool{
	var user User;
	db.Where("telephone = ? ",telephone).First(&user)
	if user.ID != 0 {
		return true;
	}

	db.AutoMigrate(&User{})
	return false
}