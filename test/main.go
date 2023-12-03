package main

import (
	"GoChatCraft/models"
	"crypto/tls"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/smtp"
	"testing"

	"github.com/jordan-wright/email"
)

func main() {
	//TestSendMail(&testing.T{})
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", "root",
		"root@123321", "127.0.0.1", 3306, "chatcraft")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&models.Relation{})
	if err != nil {
		panic(err)
	}
}
func TestSendMail(t *testing.T) {
	e := email.NewEmail()

	mailUserName := "1929509811@qq.com" //邮箱账号
	mailPassword := "erxdycnbcbjrbcfg"  //邮箱授权码
	code := "123456"                    //发送的验证码
	Subject := "验证码发送测试"                //发送的主题

	e.From = "Get <1929509811@qq.com>"
	e.To = []string{"3265804672@qq.com"}
	e.Subject = Subject
	e.HTML = []byte("你的验证码为：<h1>" + code + "</h1>")
	err := e.SendWithTLS("smtp.qq.com:465", smtp.PlainAuth("", mailUserName, mailPassword, "smtp.qq.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.qq.com"})
	if err != nil {
		t.Fatal(err)
	}
}

//func TestSendMailQQ(t *testing.T) {
//
//	mailUserName := "whm2416@qq.com"    //邮箱账号
//	mailPassword := define.MailPassword //邮箱授权码
//	addr := "smtp.qq.com:465"           //TLS地址
//	host := "smtp.qq.com"               //邮件服务器地址
//	code := "12345678"                  //发送的验证码
//	Subject := "验证码发送测试"                //发送的主题
//
//	e := email.NewEmail()
//	e.From = "Get <whm2416@qq.com>"
//	e.To = []string{"228654416@qq.com"}
//	e.Subject = Subject
//	e.HTML = []byte("你的验证码为：<h1>" + code + "</h1>")
//	err := e.SendWithTLS(addr, smtp.PlainAuth("", mailUserName, mailPassword, host),
//		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.qq.com"})
//	if err != nil {
//		t.Fatal(err)
//	}
//}
