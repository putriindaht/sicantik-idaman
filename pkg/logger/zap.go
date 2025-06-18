package logger

import (
	"fmt"

	"go.uber.org/zap"
)

var Log *zap.Logger

func InitLoggerZap() {
	var err error
	Log, err = zap.NewProduction()

	if err != nil {
		fmt.Println("Error init zap logger")
	}

	defer Log.Sync()
}
