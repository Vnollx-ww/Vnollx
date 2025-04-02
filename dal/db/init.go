package db

import (
	"Vnollx/pkg/viper"
	"Vnollx/pkg/zaplog"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	config   = viper.Init("sql")
	host     = config.Viper.GetString("mysql.source.host")
	port     = config.Viper.GetString("mysql.source.port")
	database = config.Viper.GetString("mysql.source.database")
	username = config.Viper.GetString("mysql.source.username")
	password = config.Viper.GetString("mysql.source.password")
	charset  = config.Viper.GetString("mysql.source.charset")
	logger   = zaplog.GetLogger()
)
var DB *gorm.DB

func init() {
	args := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)
	db, err := gorm.Open(mysql.Open(args), &gorm.Config{})
	if err != nil {
		logger.Panic("failed to connect to database", zap.Error(err))
	}
	//迁移
	_ = db.AutoMigrate(&User{})
	_ = db.AutoMigrate(&Folder{})
	_ = db.AutoMigrate(&Share{})
	_ = db.AutoMigrate(&File{})
	_ = db.AutoMigrate(&Notification{})
	//_ = db.AutoMigrate(&FolderRelation{})
	//初始化后的db赋给全局变量DB，没赋值的化DB就是一个nil
	DB = db
}
