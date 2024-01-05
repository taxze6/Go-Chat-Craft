package models

import (
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID       uint           `gorm:"primaryKey"`
	CreateAt time.Time      `gorm:"autoCreateTime"`
	UpdateAt time.Time      `gorm:"autoUpdateTime"`
	DeleteAt gorm.DeletedAt `gorm:"index"`
}

type UserBasic struct {
	Model
	Name          string
	PassWord      string
	Avatar        string
	Gender        string `gorm:"column:gender;default:male;type:varchar(6)"`
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9}$)"`
	Email         string `valid:"email"`
	Motto         string
	Identity      string
	ClientIp      string `valid:"ipv4"`
	ClientPort    string
	Salt          string
	LoginTime     *time.Time `gorm:"column:login_time"`
	HeartBeatTime *time.Time `gorm:"column:heart_beat_time"`
	LoginOutTime  *time.Time `gorm:"column:login_out_time"`
	IsLoginOut    bool
	DeviceInfo    string
}

type UserResponse struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Avatar     string `json:"avatar"`
	Motto      string `json:"motto"`
	Identity   string `json:"identity"`
	ClientIp   string `json:"client_ip"`
	ClientPort string `json:"client_port"`
}

func (table *UserBasic) UserTableName() string {
	return "user_basic"
}
