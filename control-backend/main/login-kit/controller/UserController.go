package controller

import (
	"control-backend/login-kit/common"
	"control-backend/login-kit/model"
	"control-backend/login-kit/response"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

func Register(ctx *gin.Context) {
	db := common.GetDB()
	//获取参数
	name := ctx.PostForm("name")
	uid := ctx.PostForm("uid")
	password := ctx.PostForm("password")
	//数据验证

	if uid == "" && password == "" && name == "" {
		json := make(map[string]interface{})
		ctx.BindJSON(&json)
		uid = json["uid"].(string)
		password = json["password"].(string)
		name = json["name"].(string)
	}

	if len(uid) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity, 442, nil, "Id必须11位")
		return
	}
	if len(password) < 8 {
		response.Response(ctx, http.StatusUnprocessableEntity, 442, nil, "密码必须大于8位")
		return
	}
	log.Println(name, uid, password)

	//判断Id是否存在
	if isTelephoneExist(db, uid) {
		response.Response(ctx, http.StatusUnprocessableEntity, 442, nil, "Id已注册")
		return
	}

	//创建用户
	hashdPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) //密码hash化
	if err != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 500, nil, "加密错误")
		return
	}
	newUser := model.User{
		Name:     name,
		Uid:      uid,
		Password: string(hashdPassword),
	}
	if err := db.Create(&newUser).Error; err != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 442, nil, err.Error())
		return
	}

	response.Success(ctx, nil, "注册成功")

}

func isTelephoneExist(db *gorm.DB, uid string) bool {
	var user model.User
	db.Where("uid=?", uid).First(&user)
	return user.ID != 0
}

func Login(ctx *gin.Context) {
	db := common.GetDB()
	//获取参数
	uid := ctx.PostForm("uid")
	password := ctx.PostForm("password")

	if uid == "" && password == "" {
		json := make(map[string]interface{})
		ctx.BindJSON(&json)
		uid = json["uid"].(string)
		password = json["password"].(string)
	}

	//数据校验
	if len(uid) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity, 442, nil, "Id必须11位")
		return
	}
	if len(password) < 8 {
		response.Response(ctx, http.StatusUnprocessableEntity, 442, nil, "密码必须大于8位")
		return
	}

	//判断Id是否存在
	var user model.User
	db.Where("uid=?", uid).First(&user)
	if user.ID == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 442, nil, "Id不存在")
		return
	}

	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 400, nil, "密码错误")
		return
	}

	//发送token
	token, err := common.GetToken(user)
	if err != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 500, nil, "系统token获取异常")
		return
	}

	response.Success(ctx, gin.H{"token": token}, "登录成功")

}

// 实现更改密码操作
func AlterPasswd(ctx *gin.Context) {
	db := common.GetDB()
	//获取参数
	uid := ctx.PostForm("uid")
	password := ctx.PostForm("password")
	newpasswd := ctx.PostForm("newpassword")
	newuid := ctx.PostForm("newuid")
	//如果form-data数据行不通，试试json格式
	if uid == "" && password == "" && newpasswd == "" && newuid == "" {
		json := make(map[string]interface{})
		ctx.BindJSON(&json)
		uid = json["uid"].(string)
		password = json["password"].(string)
		newpasswd = json["newpassword"].(string)
		newuid = json["newuid"].(string)
	}
	//对于未传入的参数不做修改
	if newuid == "" {
		newuid = uid
	}
	if newpasswd == "" {
		newpasswd = password
	}

	hashdPassword, err := bcrypt.GenerateFromPassword([]byte(newpasswd), bcrypt.DefaultCost) //密码hash化
	if err != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 500, nil, "新密码加密错误")
		return
	}

	//数据校验
	if len(uid) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity, 442, nil, "Id必须11位")
		return
	}
	if len(newuid) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity, 442, nil, "Id必须11位")
		return
	}
	if len(newpasswd) < 8 {
		response.Response(ctx, http.StatusUnprocessableEntity, 442, nil, "密码必须大于8位")
		return
	}

	if len(password) < 8 {
		response.Response(ctx, http.StatusUnprocessableEntity, 442, nil, "密码必须大于8位")
		return
	}

	//判断Id是否存在
	var user model.User
	db.Where("uid=?", uid).First(&user)
	if user.ID == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 442, nil, "Id不存在")
		return
	}

	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.Response(ctx, http.StatusUnprocessableEntity, 400, nil, "密码错误")
		return
	}
	//走到这里就可以更改密码
	if newpasswd != "" {
		user.Password = string(hashdPassword)
	}
	if newuid != "" {
		user.Uid = string(newuid)
	}

	db.Save(&user)
	response.Success(ctx, nil, "修改成功")
}

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	response.Success(ctx, gin.H{"user": user}, "")
}
