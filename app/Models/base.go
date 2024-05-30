package Models

import (
	"fmt"
	"github.com/iceking2nd/rustdesk-api-server/global"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"strings"
)

func Init() *gorm.DB {
    dbType := viper.GetString("Database.Type")
    log.Printf("Database Type: %s", dbType)

    if dbType == "" {
        log.Fatalf("Database.Type is not set or empty")
        return nil
    }

    var db *gorm.DB
    var err error

    switch dbType {
    case "mysql":
        dsn := mysqlDsn(viper.GetViper())
        db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
            DisableForeignKeyConstraintWhenMigrating: true,
        })
        if err != nil {
            log.Fatalf("Failed to initialize MySQL: %v", err)
            return nil
        }
    case "sqlite":
        sqlitePath := viper.GetString("SQLite.Path")
        db, err = gorm.Open(sqlite.Open(sqlitePath), &gorm.Config{})
        if err != nil {
            log.Fatalf("Failed to initialize SQLite: %v", err)
            return nil
        }

        if db.Exec("PRAGMA journal_mode = WAL;").Error != nil {
            log.Printf("Warning: Failed to set SQLite journal mode to WAL, error: %v", err)
        } else {
            log.Println("SQLite journal mode set to WAL successfully.")
        }

        var journalMode string
        if db.Raw("PRAGMA journal_mode;").Scan(&journalMode).Error == nil {
            log.Printf("Current SQLite journal mode: %s", journalMode)
            if strings.ToUpper(journalMode) != "WAL" {
                log.Println("Warning: Journal mode is not set to WAL despite the attempt.")
            }
        }
    default:
        log.Fatalf("Unsupported database type: %s", dbType)
        return nil
    }

    if db == nil {
        log.Fatalf("Failed to initialize database: %v", err)
        return nil
    }

    dbSQL, err := db.DB()
    if err != nil {
        log.Fatalf("Failed to get raw DB instance: %v", err)
        return nil
    }

    dbSQL.SetMaxIdleConns(5)
    dbSQL.SetMaxOpenConns(20)

    if global.LogLevel >= 5 {
        return db.Debug()
    }
    return db
}

func mysqlDsn(v *viper.Viper) string {
	unixSocketPath := v.GetString("MySQL.UnixSocket")
	if unixSocketPath != "" {
		return fmt.Sprintf(
			"%s:%s@unix(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			v.GetString("MySQL.User"),
			v.GetString("MySQL.Pass"),
			unixSocketPath,
			v.GetString("MySQL.DB"),
		)
	} else {
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			v.GetString("MySQL.User"),
			v.GetString("MySQL.Pass"),
			v.GetString("MySQL.Host"),
			v.GetInt("MySQL.Port"),
			v.GetString("MySQL.DB"),
		)
	}
}
