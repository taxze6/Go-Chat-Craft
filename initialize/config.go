package initialize

import (
	"GoChatCraft/global"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func InitConfig() {
	//Instantiate an object
	v := viper.New()

	configFile := "../GoChatCraft/config-release.yaml"

	//Read configuration file
	v.SetConfigFile(configFile)

	//Read file
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	//Put the data into global.ServerConfig.
	if err := v.Unmarshal(&global.ServiceConfig); err != nil {
		panic(err)
	}

	zap.S().Info("Configuration information.", global.ServiceConfig)
}
