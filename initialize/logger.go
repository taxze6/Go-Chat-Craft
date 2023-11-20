package initialize

import (
	"go.uber.org/zap"
	"log"
)

func InitLogger() {
	//Initialize logging.
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("The initialization of the logging failed.", err.Error())
	}
	//Use global logger.
	zap.ReplaceGlobals(logger)
}
