package service

import (
	"GoChatCraft/dao"
	"GoChatCraft/models"
)

func GetEmailCode(mailTo string) error {
	option := &models.MailOptions{
		MailHost: "smtp.qq.com",
		MailPort: 465,
		MailUser: "1929509811@qq.com",
		MailPass: "erxdycnbcbjrbcfg",
		MailTo:   mailTo,
		Subject:  "email code test",
		Body:     "",
	}
	err := dao.SendMailCode(option, "register", 5)
	if err != nil {
		return err
	}
	return nil
}

func CheckEmailCode(email string, code string) error {
	err := dao.ValidateMailCode(email, code, "register")
	if err != nil {
		return err
	}
	return nil
}
