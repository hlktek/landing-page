package logger

import (
	"oauth2-go-service/config"

	logrus "github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02T15:04:05.000Z07:00"
	customFormatter.FullTimestamp = true
	customFormatter.ForceColors = true
	Logger.SetFormatter(customFormatter)
	// open a file
	// f, err := os.OpenFile(config.GetConfig("log_path").(string), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	// if err != nil {
	// 	fmt.Printf("error opening file: %v", err)
	// }
	Logger.SetNoLock()
	// Logger.SetOutput(f)
	logLevel := config.GetConfig("LOG_LEVEL")

	if logLevel == "" {
		Logger.SetLevel(logrus.InfoLevel)
	} else {
		if logLevel == "DEBUG" {
			Logger.SetLevel(logrus.DebugLevel)
		} else if logLevel == "INFO" {
			Logger.SetLevel(logrus.InfoLevel)
		} else if logLevel == "WARN" {
			Logger.SetLevel(logrus.WarnLevel)
		} else if logLevel == "ERROR" {
			Logger.SetLevel(logrus.ErrorLevel)
		}
	}
}

//Info log info
func Info(fields logrus.Fields, message string, params ...interface{}) {
	Logger.WithFields(fields).Infof(message, params...)
}

//Debug log debug
func Debug(fields logrus.Fields, message string, params ...interface{}) {
	Logger.WithFields(fields).Debugf(message, params...)
}

//Error log error
func Error(fields logrus.Fields, message string, params ...interface{}) {
	Logger.WithFields(fields).Errorf(message, params...)
}

// func Infofs(fields logrus.Fields, message string, params ...interface{}) {
// 	Logger.WithFields(fields).Infof(message, params)
// }
// func Warn(key string, message string) {
// 	Logger.WithField("transaction_id", key).Warn(message)
// }
// func Warnf(key string, message string, params ...interface{}) {
// 	Logger.WithField("transaction_id", key).Warnf(message, params...)
// }
// func Warnfs(fields logrus.Fields, message string, params ...interface{}) {
// 	Logger.WithFields(fields).Warnf(message, params...)
// }
