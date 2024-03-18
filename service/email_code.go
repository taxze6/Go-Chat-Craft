package service

import (
	"GoChatCraft/dao"
	"GoChatCraft/models"
)

func GetEmailCode(mailTo string, from string) error {
	option := &models.MailOptions{
		MailHost: "smtp.qq.com",
		MailPort: 465,
		MailUser: "1929509811@qq.com",
		MailPass: "aqvgewdgxhygdecf",
		MailTo:   mailTo,
		Subject:  "chat craft email code test",
		Body:     "",
	}
	err := dao.SendMailCode(option, from, 1)
	if err != nil {
		return err
	}
	return nil
}

func CheckEmailCode(email string, code string, from string) error {
	err := dao.ValidateMailCode(email, code, from)
	if err != nil {
		return err
	}
	return nil
}
