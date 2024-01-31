package config

// MysqlConfig mysql信息配置
type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Name     string `mapstructure:"name" json:"Name"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
}

type RedisConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
}

type RabbitMQConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     string `mapstructure:"port" json:"port"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
}

// 对应yaml文件结构
type ServiceConfig struct {
	Port           int            `mapstructure:"port" json:"port"`
	Host           string         `mapstructure:"host" json:"host"`
	DB             MysqlConfig    `mapstructure:"mysql" json:"mysql"`
	RedisDB        RedisConfig    `mapstructure:"redis" json:"redis"`
	RabbitMQConfig RabbitMQConfig `mapstructure:"rabbitmq" json:"rabbitmq"`
}
