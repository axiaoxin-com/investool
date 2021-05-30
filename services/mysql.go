package services

import (
	"context"
	"time"

	"github.com/axiaoxin-com/goutils"
	"github.com/axiaoxin-com/logging"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 获取带有ctx 和自定义 logger 的 gorm db 实例
func DB(ctx context.Context) *gorm.DB {
	env := viper.GetString("env")
	if env == "unittest" {
		dbname := viper.GetString("unittest.dbname")
		if dbname == "" {
			dbname = "/tmp/pinklady_test.db"
		}
		db, err := goutils.NewGormSQLite3(goutils.DBConfig{DBName: dbname, GormConfig: &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		}})
		if err != nil {
			panic(err)
		}
		return db
	}

	db, err := goutils.GormMySQL(env)
	if err != nil {
		panic(env + " get gorm mysql instance error:" + err.Error())
	}
	logging.Debug(ctx, "using gorm mysql:"+env)
	db = db.Session(&gorm.Session{
		Logger: logging.NewGormLogger(zap.InfoLevel, zap.DebugLevel, viper.GetDuration("logging.access_logger.slow_threshold")*time.Millisecond),
	})
	return db.WithContext(ctx)
}
