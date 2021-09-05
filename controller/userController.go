package controller

import (
	"login4/common"
	"login4/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"login4/util"
	"net/http"
)

func Register(ctx *gin.Context) {
	DB := common.GetDB();

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
		name = util.RandomName(10)
	}

	//判断手机号是否存在
	if isTelephoneExist(DB,telephone) {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户存在", "number": len(telephone)})
		return
	}

	//创建用户
	newUser := model.User{
		Name: name,
		Telephone: telephone,
		Password: password,
	}
	DB.Create(&newUser);
	//返回结果
	ctx.JSON(200, gin.H{
		"message": "注册成功",
	})
}

func isTelephoneExist(db *gorm.DB,telephone string) bool{
	var user model.User;
	db.Where("telephone = ? ",telephone).First(&user)
	if user.ID != 0 {
		return true;
	}

	db.AutoMigrate(&model.User{})
	return false
}