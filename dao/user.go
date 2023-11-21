package dao

import (
	"GoChatCraft/common"
	"GoChatCraft/global"
	"GoChatCraft/models"
	"errors"
	"go.uber.org/zap"
	"strconv"
	"time"
)

func GetUserList() ([]*models.UserBasic, error) {
	var list []*models.UserBasic
	//tx.RowsAffected is a method that returns the number of affected rows, used to check if a query has returned any results.
	//Use if tx.RowsAffected == 0 to check if the query result is empty. If it is empty, it means no user records were found.
	if tx := global.DB.Find(&list); tx.RowsAffected == 0 {
		zap.S().Info("failed to retrieve user list")
		return nil, errors.New("failed to retrieve user list")
	}
	return list, nil
}

func FindUserByNameAndPwd(name string, password string) (*models.UserBasic, error) {
	user := models.UserBasic{}
	if tx := global.DB.Where("name = ? and pass_word = ?", name, password).First(&user); tx.RowsAffected == 0 {
		zap.S().Info("the user was not found")
		return nil, errors.New("the user was not found")
	}
	//Get the current timestamp and convert it to a string type.
	t := strconv.Itoa(int(time.Now().Unix()))
	//Perform MD5 encryption.
	temp := common.Md5encoder(t)
	if tx := global.DB.Model(&user).Where("id = ?", user.ID).Update("identity", temp); tx.RowsAffected == 0 {
		zap.S().Info("failed to write identity")
		return nil, errors.New("failed to write identity")
	}
	return &user, nil
}

func FindUserByName(name string) (*models.UserBasic, error) {
	user := models.UserBasic{}
	if tx := global.DB.Where("name = ?", name).First(&user); tx.RowsAffected == 0 {
		zap.S().Info("couldn't find any information about this user")
		return nil, errors.New("couldn't find any information about this user")
	}
	return &user, nil
}

func FindUserByNameWithRegister(name string) (*models.UserBasic, error) {
	user := models.UserBasic{}
	if tx := global.DB.Where("name = ?", name).First(&user); tx.RowsAffected == 1 {
		zap.S().Info("the current username already exists")
		return nil, errors.New("the current username already exists")
	}
	return &user, nil
}

func FindUserId(ID uint) (*models.UserBasic, error) {
	user := models.UserBasic{}
	if tx := global.DB.Where(ID).First(&user); tx.RowsAffected == 0 {
		zap.S().Info("the user was not found")
		return nil, errors.New("the user was not found")
	}
	return &user, nil
}

func CreateUser(user models.UserBasic) (*models.UserBasic, error) {
	tx := global.DB.Create(&user)
	if tx.RowsAffected == 0 {
		zap.S().Info("failed to add a new user")
		return nil, errors.New("failed to add a new user")
	}
	return &user, nil
}

func UpdateUser(user models.UserBasic) (*models.UserBasic, error) {
	tx := global.DB.Model(&user).Updates(models.UserBasic{
		Name:     user.Name,
		PassWord: user.PassWord,
		Avatar:   user.Avatar,
		Gender:   user.Gender,
		Phone:    user.Phone,
		Email:    user.Email,
		Salt:     user.Salt,
	})
	if tx.RowsAffected == 0 {
		zap.S().Info("failed to update the user")
		return nil, errors.New("failed to update the user")
	}
	return &user, nil
}

func DeleteUser(user models.UserBasic) error {
	if tx := global.DB.Delete(&user); tx.RowsAffected == 0 {
		zap.S().Info("failed to delete the user")
		return errors.New("failed to delete the user")
	}
	return nil
}
