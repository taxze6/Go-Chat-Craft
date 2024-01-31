package service

import (
	"GoChatCraft/common"
	"GoChatCraft/dao"
	"GoChatCraft/global"
	"GoChatCraft/middlewear"
	"GoChatCraft/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// GetUserList Get All User List
//
//	@Summary		List Get User List
//	@Description	User List
//	@Tags			test
//	@Accept			json
//	@Router			/user/list [get]
func GetUserList(ctx *gin.Context) {
	list, err := dao.GetUserList()
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    -1,
			"message": "Failed to retrieve the user list.",
		})
		return
	}
	ctx.JSON(http.StatusOK, list)
}

func LoginByNameAndPassWord(ctx *gin.Context) {
	getData, _ := ctx.GetRawData()
	var body map[string]string
	_ = json.Unmarshal(getData, &body)
	name := body["name"]
	encryptionPassword := body["password"]
	//First, check if the username exists, then proceed to check the password.
	data, err := dao.FindUserByName(name)
	if err != nil {
		common.RespFail(ctx.Writer, "Login Failed", "Login Failed.")
		return
	}
	if data.Name == "" {
		common.RespFail(ctx.Writer, "Username Does Not Exist.", "Username Does Not Exist.")
		return
	}
	password, err := common.RsaDecoder(encryptionPassword)
	if err != nil {
		zap.S().Info("RSA Decryption Failed")
		common.RespFail(ctx.Writer, "Password parsing error!", "Password parsing error!")
		return
	}
	ok := common.CheckPassWord(password, data.Salt, data.PassWord)
	if !ok {
		common.RespFail(ctx.Writer, "Incorrect Password.", "Incorrect Password.")
		return
	}
	userInfo, err := dao.FindUserByNameAndPwd(name, data.PassWord)
	if err != nil {
		zap.S().Info("Login Failed.", err)
		common.RespFail(ctx.Writer, "Login Failed.", "Login Failed.")
		return
	}
	//Using JWT for authentication.
	token, err := middlewear.GenerateToken(userInfo.ID, "cc")
	if err != nil {
		zap.S().Info("Failed to Generate Token", err)
		common.RespFail(ctx.Writer, "Failed to Generate Token.", "Failed to Generate Token.")
		return
	}
	response := models.UserResponse{
		ID:         userInfo.ID,
		Name:       userInfo.Name,
		Email:      userInfo.Email,
		Phone:      userInfo.Phone,
		Avatar:     userInfo.Avatar,
		Motto:      userInfo.Motto,
		Identity:   userInfo.Identity,
		ClientIp:   userInfo.ClientIp,
		ClientPort: userInfo.ClientPort,
	}
	common.RespOk(ctx.Writer, gin.H{
		"token": token,
		"user":  response}, "Login Successful.")
}

func NewUser(ctx *gin.Context) {
	user := models.UserBasic{}
	getData, _ := ctx.GetRawData()
	var body map[string]string
	_ = json.Unmarshal(getData, &body)
	user.Name = body["name"]
	user.Email = body["email"]
	encryptionPassword := body["password"]
	encryptionRePassword := body["repassword"]
	password, err := common.RsaDecoder(encryptionPassword)
	if err != nil {
		zap.S().Info("RSA Decryption Failed")
		common.RespFail(ctx.Writer, "Password parsing error!", "Password parsing error!")
		return
	}
	repassword, err := common.RsaDecoder(encryptionRePassword)
	if err != nil {
		zap.S().Info("RSA Decryption Failed")
		common.RespFail(ctx.Writer, "RePassword parsing error!", "RePassword parsing error!")
		return
	}
	if user.Name == "" || password == "" || repassword == "" {
		common.RespFail(ctx.Writer, "Username or password cannot be empty!", "Username or password cannot be empty!")
		return
	}
	if user.Email == "" {
		common.RespFail(ctx.Writer, "Email cannot be empty!", "Email cannot be empty!")
		return
	}
	//查询用户是否存在
	_, err = dao.FindUserByNameWithRegister(user.Name)
	if err != nil {
		common.RespFail(ctx.Writer, "The user has already registered!", "The user has already registered!")
		return
	}
	//查询邮箱是否已被注册
	_, err = dao.FindUserByEmailWithRegister(user.Name)
	if err != nil {
		common.RespFail(ctx.Writer, "The email has already registered!", "The email has already registered!")
		return
	}

	if password != repassword {
		common.RespFail(ctx.Writer, "The passwords entered do not match!", "The passwords entered do not match!")
		return
	}
	err = GetEmailCode(user.Email, global.Register)
	if err != nil {
		zap.S().Info("failed to send verification code")
		common.RespFail(ctx.Writer, "failed to send verification code!", "failed to send verification code!")
		return
	}
	common.RespOk(ctx.Writer, "Verification code sent successfully.", "Verification code sent successfully.")
}

func CheckRegisterEmailCode(ctx *gin.Context) {
	user := models.UserBasic{}
	getData, _ := ctx.GetRawData()
	var body map[string]string
	_ = json.Unmarshal(getData, &body)
	user.Name = body["name"]
	user.Email = body["email"]
	password := body["password"]
	ps, err := common.RsaDecoder(password)
	code := body["code"]
	err = CheckEmailCode(user.Email, code, global.Register)
	if err != nil {
		zap.S().Info("incorrect verification code")
		common.RespFail(ctx.Writer, "Incorrect verification code!", "Incorrect verification code!")
		return
	}
	salt := fmt.Sprintf("%d", rand.Int31())
	//加密密码
	user.PassWord = common.SaltPassWord(ps, salt)
	user.Salt = salt
	t := time.Now()
	user.LoginTime = &t
	user.LoginOutTime = &t
	user.HeartBeatTime = &t
	_, err = dao.CreateUser(user)
	if err != nil {
		return
	}
	info, _ := dao.FindUserByName(user.Name)
	//Using JWT for authentication.
	token, err := middlewear.GenerateToken(info.ID, "cc")
	if err != nil {
		zap.S().Info("Failed to Generate Token", err)
		common.RespFail(ctx.Writer, "Failed to Generate Token", "Failed to Generate Token")
		return
	}
	userInfo := models.UserResponse{
		ID:         info.ID,
		Name:       info.Name,
		Email:      info.Email,
		Phone:      info.Phone,
		Avatar:     info.Avatar,
		Motto:      info.Motto,
		Identity:   info.Identity,
		ClientIp:   info.ClientIp,
		ClientPort: info.ClientPort,
	}
	common.RespOk(ctx.Writer, gin.H{
		"token": token,
		"user":  userInfo,
	}, "New user added successfully！")
}

func EmailLogin(ctx *gin.Context) {
	getData, _ := ctx.GetRawData()
	var body map[string]string
	_ = json.Unmarshal(getData, &body)
	email := body["email"]
	_, err := dao.FindUserByEmailWithLogin(email)
	if err != nil {
		common.RespFail(ctx.Writer, "Couldn't find any information about this email.", "Couldn't find any information about this email.")
		return
	}
	err = GetEmailCode(email, global.LoginEmail)
	if err != nil {
		zap.S().Info("failed to send verification code")
		common.RespFail(ctx.Writer, "failed to send verification code!", "failed to send verification code!")
		return
	}
	common.RespOk(ctx.Writer, "Verification code sent successfully.", "Verification code sent successfully.")
}

func CheckLoginEmailCode(ctx *gin.Context) {
	getData, _ := ctx.GetRawData()
	var body map[string]string
	_ = json.Unmarshal(getData, &body)
	email := body["email"]
	code := body["code"]
	err := CheckEmailCode(email, code, global.LoginEmail)
	if err != nil {
		zap.S().Info("incorrect verification code")
		common.RespFail(ctx.Writer, "Incorrect verification code!", "Incorrect verification code!")
		return
	}
	//查询用户数据
	user, err := dao.FindUserByEmailWithLogin(email)
	if err != nil {
		common.RespFail(ctx.Writer, "Couldn't find any information about this email.", "Couldn't find any information about this email.")
		return
	}
	t := time.Now()
	user.LoginTime = &t
	user.LoginOutTime = &t
	user.HeartBeatTime = &t
	//Using JWT for authentication.
	token, err := middlewear.GenerateToken(user.ID, "cc")
	if err != nil {
		zap.S().Info("Failed to Generate Token", err)
		common.RespFail(ctx.Writer, "Failed to Generate Token", "Failed to Generate Token")
		return
	}
	userInfo := models.UserResponse{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		Phone:      user.Phone,
		Avatar:     user.Avatar,
		Motto:      user.Motto,
		Identity:   user.Identity,
		ClientIp:   user.ClientIp,
		ClientPort: user.ClientPort,
	}
	common.RespOk(ctx.Writer, gin.H{
		"token": token,
		"user":  userInfo,
	}, "Login Successful.")
}

func FindUserWithUserName(ctx *gin.Context) {
	getData, _ := ctx.GetRawData()
	var body map[string]string
	_ = json.Unmarshal(getData, &body)
	name := body["name"]
	userInfo, err := dao.FindUserByName(name)
	if err != nil {
		common.RespFail(ctx.Writer, "couldn't find any information about this user", "couldn't find any information about this user")
		return
	}
	reUserInfo := models.UserResponse{
		ID:         userInfo.ID,
		Name:       userInfo.Name,
		Email:      userInfo.Email,
		Phone:      userInfo.Phone,
		Avatar:     userInfo.Avatar,
		Motto:      userInfo.Motto,
		Identity:   userInfo.Identity,
		ClientIp:   userInfo.ClientIp,
		ClientPort: userInfo.ClientPort,
	}
	common.RespOk(ctx.Writer, reUserInfo, "The user has been found.")
}

func UpdateUser(ctx *gin.Context) {
	userId, _ := strconv.Atoi(ctx.GetHeader("UserId"))
	getData, _ := ctx.GetRawData()
	var body map[string]string
	_ = json.Unmarshal(getData, &body)
	avatar := body["avatar"]
	motto := body["motto"]
	name := body["name"]
	phone := body["phone"]
	email := body["email"]
	userInfo, err := dao.FindUserId(uint(userId))
	if userInfo.Avatar != avatar {
		userInfo.Avatar = avatar
	}

	if userInfo.Motto != motto {
		userInfo.Motto = motto
	}

	if userInfo.Name != name {
		data, _ := dao.FindUserByName(name)
		if data != nil {
			common.RespFail(ctx.Writer, "The user name already exists!", "The user name already exists!")
			return
		}
		userInfo.Name = name
	}

	if userInfo.Phone != phone {
		userInfo.Phone = phone
	}

	if userInfo.Email != email {
		userInfo.Email = email
	}
	newUserInfo, err := dao.UpdateUser(*userInfo)
	if err != nil {
		common.RespFail(ctx.Writer, "Fail to modify.", "Fail to modify.")
		return
	}

	common.RespOk(ctx.Writer, newUserInfo, "The account information is successfully modified!")
}

func UpdateUserPassword(ctx *gin.Context) {
	userId, _ := strconv.Atoi(ctx.GetHeader("UserId"))
	getData, _ := ctx.GetRawData()
	var body map[string]string
	_ = json.Unmarshal(getData, &body)
	encryptionPassword := body["password"]
	encryptionNewPassword := body["newPassword"]
	password, err := common.RsaDecoder(encryptionPassword)
	newPassword, err := common.RsaDecoder(encryptionNewPassword)
	if err != nil {
		zap.S().Info("Cryptographic error")
		common.RespFail(ctx.Writer, "Cryptographic error!", "Cryptographic error!")
		return
	}
	userInfo, err := dao.FindUserId(uint(userId))
	ok := common.CheckPassWord(password, userInfo.Salt, userInfo.PassWord)
	if !ok {
		common.RespFail(ctx.Writer, "The old password is incorrect.", "The old password is incorrect.")
		return
	}
	userInfo.PassWord = common.SaltPassWord(newPassword, userInfo.Salt)
	newUserInfo, err := dao.UpdateUser(*userInfo)
	if err != nil {
		common.RespFail(ctx.Writer, "Fail to modify.", "Fail to modify.")
		return
	}

	common.RespOk(ctx.Writer, newUserInfo, "The account information is successfully modified!")
}
