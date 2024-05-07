package logging

import (
	"dcss/global"
	"fmt"
	rotateLogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"time"
)

func Init() error {
	logPath := viper.GetString("LOG.PATH")
	logAbsDir := filepath.Join(global.BaseDir, filepath.Dir(logPath))
	logName := filepath.Base(logPath)
	logMaxAge := viper.GetInt("Log.MaxAge")

	if _, err := os.Stat(logAbsDir); os.IsNotExist(err) {
		err := os.MkdirAll(logAbsDir, 0o755)
		if err != nil {
			return fmt.Errorf("create log dir failed, err: %s", err.Error())
		}
	}

	logDateName := fmt.Sprintf("%s/%%Y%%m%%d-%s", logAbsDir, logName)
	logLinkName := fmt.Sprintf("%s/%s", logAbsDir, logName)

	logf, err := rotateLogs.New(
		logDateName,
		rotateLogs.WithLinkName(logLinkName),
		rotateLogs.WithRotationTime(24*time.Hour),
		rotateLogs.WithMaxAge(time.Duration(logMaxAge)*24*time.Hour),
	)

	if err != nil {
		return fmt.Errorf("new ratatelog failed, err: %s", err.Error())
	}

	global.LOG = logrus.New()
	global.LOG.SetOutput(logf)

	customLogFormatter := logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	}

	global.LOG.SetFormatter(&customLogFormatter)

	// 0： panic, 1: fatal, 2: error, 3: warn, 4: info, 5: debug, 6: trace
	global.LOG.SetLevel(logrus.Level(viper.GetInt("LOG.LEVEL")))
	global.LOG.SetReportCaller(true) //  Logging Method Name

	return nil
}
