package logger

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"

	"github.com/bdzhalalov/resilient-scatter-gather/config"
)

func Logger(config *config.Config) *logrus.Logger {

	log := logrus.New()

	level, err := logrus.ParseLevel(config.LogLevel)
	if err != nil {
		fmt.Println(err.Error())
		log.Error("Can't parse log level, setting level to DEBUG")
		level = logrus.DebugLevel
	}

	log.SetLevel(level)

	log.SetOutput(os.Stdout)

	file, err := os.OpenFile("./logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)

		log.AddHook(&writer.Hook{
			Writer: file,
			LogLevels: []logrus.Level{
				logrus.PanicLevel,
				logrus.FatalLevel,
				logrus.ErrorLevel,
				logrus.WarnLevel,
				logrus.DebugLevel,
			},
		})

		log.AddHook(&writer.Hook{
			Writer: os.Stdout,
			LogLevels: []logrus.Level{
				logrus.InfoLevel,
			},
		})
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	return log
}
