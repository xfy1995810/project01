package db

import (
	"dcss/global"
	"dcss/models"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	rotateLogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB 初始化数据库
func InitDB() error {
	dbPath := path.Join(global.BaseDir, viper.GetString("DB.PATH"))
	var err error

	dbLogger, err := createDbLog()
	if err != nil {
		return err
	}

	global.DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: dbLogger.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect database, err:" + err.Error())
	}

	sqlDB, err := global.DB.DB()
	if err != nil {
		return fmt.Errorf("get *sql.DB error, err:" + err.Error())
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	err = migrateDB()
	if err != nil {
		return err
	}
	return nil
}

func createDbLog() (logger.Interface, error) {
	logPath := viper.GetString("LOG.PATH")
	logAbsDir := filepath.Join(global.BaseDir, filepath.Dir(logPath))
	logName := "db.log"
	logMaxAge := viper.GetInt("Log.MaxAge")

	if _, err := os.Stat(logAbsDir); os.IsNotExist(err) {
		err := os.MkdirAll(logAbsDir, 0o755)
		if err != nil {
			return nil, fmt.Errorf("create log dir failed, err:%s\n", err.Error())
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
		return nil, fmt.Errorf("new ratatelog failed, err: %s\n", err.Error())
	}

	newLogger := logger.New(
		log.New(logf, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful:                  false,         // Disable color
		},
	)

	return newLogger, nil
}

func migrateDB() error {
	err := global.DB.AutoMigrate(
		&models.SysUser{},
		&models.SysRole{},
		&models.SysApi{},
		&models.SysMenu{},
		&models.SampleInfos{},
	)
	if err != nil {
		return fmt.Errorf("auto migrate db failed, err: %s", err)
	}

	// 初始化超级管理员用户
	err = initUserData()
	if err != nil {
		global.LOG.Errorln("初始化超级管理员数据失败，err: ", err)
		return err
	}

	// 初始化Api
	err = initApiData()
	if err != nil {
		global.LOG.Errorln("初始化Api数据失败，err: ", err)
		return err
	}

	// 初始化菜单数据
	err = initMenuData()
	if err != nil {
		global.LOG.Errorln("初始化菜单数据失败，err: ", err)
		return err
	}

	return nil
}
