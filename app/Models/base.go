package Models

import (
	"fmt"
	"github.com/iceking2nd/rustdesk-api-server/global"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func Init() *gorm.DB {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("MySQL.User"),
		viper.GetString("MySQL.Pass"),
		viper.GetString("MySQL.Host"),
		viper.GetInt("MySQL.Port"),
		viper.GetString("MySQL.DB"),
	)

	if db, errOpenDB := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	}); errOpenDB != nil {
		log.Fatalln(errOpenDB.Error())
	} else {
		if dbSQL, errGetDBSQL := db.DB(); errGetDBSQL != nil {
			log.Fatalln(errGetDBSQL.Error())
		} else {
			dbSQL.SetMaxIdleConns(5)
			dbSQL.SetMaxOpenConns(20)
			if global.LogLevel >= 5 {
				return db.Debug()
			}
			return db
		}
	}
	return nil
}
