package logger

import (
	"fmt"
	"go.uber.org/zap"
)

var logger *zap.Logger
var err error

func init() {
	logger, err = zap.NewProduction()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func GetLogger() *zap.Logger {
	return logger
}
