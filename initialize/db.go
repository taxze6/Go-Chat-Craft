package initialize

import (
	"GoChatCraft/global"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func InitDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", global.ServiceConfig.DB.User,
		global.ServiceConfig.DB.Password, global.ServiceConfig.DB.Host, global.ServiceConfig.DB.Port, global.ServiceConfig.DB.Name)
	newLogger := logger.New(log.New(os.Stdout, "\r\n", log.Lshortfile),
		logger.Config{
			SlowThreshold:             time.Second, //慢Sql阈值（当执行的 SQL 查询或操作的执行时间超过该阈值时，会被认为是慢 SQL。 慢 SQL是指执行时间较长的 SQL 查询或操作，可能会对系统性能产生影响。通过设置慢 SQL 的阈值，可以在日志中标记并记录执行时间超过阈值的 SQL 查询或操作，以便进行性能分析和优化。）
			LogLevel:                  logger.Info, //日志级别（logger.Info表示只输出信息级别及以上的日志）
			IgnoreRecordNotFoundError: true,        //忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,        //彩色打印
		},
	)
	var err error
	//将获取到的连接赋值到global.DB
	global.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
}

func InitRedis() {
	opt := redis.Options{
		Addr:     fmt.Sprintf("%s:%d", global.ServiceConfig.RedisDB.Host, global.ServiceConfig.RedisDB.Port), // redis地址
		Password: "",                                                                                         // no password set
		DB:       0,                                                                                          // if you want to use default DB,set to 0
	}
	global.RedisDB = redis.NewClient(&opt)
}
