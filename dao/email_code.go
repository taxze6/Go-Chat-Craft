package dao

import (
	"GoChatCraft/global"
	"GoChatCraft/models"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jordan-wright/email"
	"go.uber.org/zap"
	"math/rand"
	"net/smtp"
	"time"
)

var ctx = context.Background()

func GenerateRandomCode(codeLen int) string {
	s := "1234567890"
	code := ""
	// import random seed, or that random code will always be the same one
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < codeLen; i++ {
		code += string(s[rand.Intn(len(s))])
	}
	return code
}

// GenMailCodeKey Where does the captcha business come from
func GenMailCodeKey(mailAddr, from string) string {
	return "MAIL-CODE-" + from + "-" + mailAddr
}

func SendMailCode(mailOption *models.MailOptions, from string, ttl int) error {
	// generate random validation code
	code := GenerateRandomCode(6)
	// 先插入redis，再发邮件
	key := GenMailCodeKey(mailOption.MailTo, from)
	// 如果code没有过期，是不允许再发送的
	success, err := global.RedisDB.SetNX(ctx, key, code, time.Duration(ttl)*time.Minute).Result()
	if err != nil {
		return err
	}
	if !success {
		return errors.New("the code already exists")
	}

	// 发邮件
	options := &models.MailOptions{
		MailHost: "smtp.qq.com",
		MailPort: 465,
		MailUser: mailOption.MailUser,
		MailPass: mailOption.MailPass,
		MailTo:   mailOption.MailTo,
		Subject:  mailOption.Subject,
		Body:     mailOption.Body,
	}
	err = MailSend(options, code)
	if err != nil {
		return err
	}
	return nil
}

func MailSend(options *models.MailOptions, code string) error {
	e := email.NewEmail()
	e.From = "Chat Craft<1929509811@qq.com>"
	e.To = []string{options.MailTo}
	e.Subject = options.Subject
	e.HTML = []byte("<h1>Your Validation Code is " + code + "</h1>")
	err := e.SendWithTLS("smtp.qq.com:465", smtp.PlainAuth("", options.MailUser, options.MailPass, "smtp.qq.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.qq.com"})
	if err != nil {
		zap.S().Info("email delivery failed: %w", err)
		return errors.New("email delivery failed")
	}
	return nil
}

func ValidateMailCode(tag, inputCode, from string) error {
	key := GenMailCodeKey(tag, from)
	code, err := global.RedisDB.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return errors.New("")
		}
		return err
	}

	// 对比后马上删除
	err = global.RedisDB.Del(ctx, key).Err()
	if err != nil {
		fmt.Printf("redis del fail %v\n", err)
		return err
	}

	if inputCode != code {
		return errors.New("验证码不对")
	}

	return nil
}
