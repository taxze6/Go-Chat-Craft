package service

import (
	"GoChatCraft/common"
	"GoChatCraft/dao"
	"GoChatCraft/middlewear"
	"GoChatCraft/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"time"
)

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
		ctx.JSON(200, gin.H{
			"code":    -1, //0 represents success, -1 represents failure
			"message": "Login Failed.",
		})
		return
	}
	if data.Name == "" {
		ctx.JSON(200, gin.H{
			"code":    -1,
			"message": "Username Does Not Exist.",
		})
		return
	}
	password, err := common.RsaDecoder(encryptionPassword)
	if err != nil {
		zap.S().Info("RSA Decryption Failed")
		ctx.JSON(200, gin.H{
			"code":    -1,
			"message": "Password parsing error!",
		})
		return
	}
	ok := common.CheckPassWord(password, data.Salt, data.PassWord)
	if !ok {
		ctx.JSON(200, gin.H{
			"code":    -1,
			"message": "Incorrect Password.",
		})
		return
	}
	userInfo, err := dao.FindUserByNameAndPwd(name, data.PassWord)
	if err != nil {

		zap.S().Info("Login Failed.", err)
		return
	}
	//Using JWT for authentication.
	token, err := middlewear.GenerateToken(userInfo.ID, "cc")
	if err != nil {
		zap.S().Info("Failed to Generate Token", err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Login Successful.",
		"tokens":  token,
		"userId":  userInfo.ID,
	})
}

func NewUser(ctx *gin.Context) {
	user := models.UserBasic{}
	getData, _ := ctx.GetRawData()
	var body map[string]string
	_ = json.Unmarshal(getData, &body)
	user.Name = body["name"]
	user.Email = body["email"]
	password := body["password"]
	repassword := body["repassword"]
	if user.Name == "" || password == "" || repassword == "" {
		ctx.JSON(200, gin.H{
			"code":    -1, //  0成功   -1失败
			"message": "Username or password cannot be empty！",
			"data":    nil,
		})
		return
	}
	if user.Email == "" {
		ctx.JSON(200, gin.H{
			"code":    -1, //  0成功   -1失败
			"message": "Email cannot be empty！",
			"data":    nil,
		})
		return
	}
	//查询用户是否存在
	_, err := dao.FindUserByNameWithRegister(user.Name)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code":    -1,
			"message": "The user has already registered!",
			"data":    nil,
		})
		return
	}

	if password != repassword {
		ctx.JSON(200, gin.H{
			"code":    -1,
			"message": "The passwords entered do not match！",
			"data":    nil,
		})
		return
	}
	salt := fmt.Sprintf("%d", rand.Int31())
	//加密密码
	user.PassWord = common.SaltPassWord(password, salt)
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
	response := models.UserResponse{
		ID:         info.ID,
		Name:       info.Name,
		Email:      info.Email,
		Phone:      info.Phone,
		Avatar:     info.Avatar,
		Motto:      info.Motto,
		ClientIp:   info.ClientIp,
		ClientPort: info.ClientPort,
	}
	ctx.JSON(200, gin.H{
		"code":    0, //  0成功   -1失败
		"message": "New user added successfully！",
		"data":    response,
	})
}
