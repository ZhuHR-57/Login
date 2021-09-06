package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"login4/common"
	"login4/model"
	"login4/util"
	"net/http"
)

/**
	用户注册
 */
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
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户存在"})
		return
	}

	//创建用户
	//密码加密
	hasedPassWord,err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost);
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "加密出错"})
		return
	}

	newUser := model.User{
		Name: name,
		Telephone: telephone,
		Password: string(hasedPassWord),
	}
	DB.Create(&newUser);

	//返回结果
	ctx.JSON(200, gin.H{
		"code": 200,
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

/**
	用户登录
 */
func Login(ctx *gin.Context){
	DB := common.GetDB();

	//获取参数
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")

	//数据验证
	if len(password) < 6 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码不能少于6位"})
		return
	}
	if len(telephone) != 11 {
		//gin.H = map[string]interface{}
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位"})
		return
	}

	//判断手机号是否存在
	var user model.User;
	DB.Where("telephone = ? ",telephone).First(&user)
	if user.ID == 0 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户不存在"})
		return
	}

	//判断密码是否接收正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password));err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "密码错误"})
		return
	}

	//发放token
	token := "11"

	//返回结果
	ctx.JSON(200, gin.H{
		"code": 200,
		"data": gin.H{"token": token},
		"message": "登录成功",
	})
}