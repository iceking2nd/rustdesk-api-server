package Database

import (
	"github.com/gin-gonic/gin"
	"github.com/iceking2nd/rustdesk-api-server/app/Models"
	"gorm.io/gorm"
)

func SetContext() gin.HandlerFunc {
	db := Models.Init()
	return func(ctx *gin.Context) {
		ctx.Set("db", db)
		ctx.Next()
	}
}

func GetDB(ctx *gin.Context) *gorm.DB {
	idb, ok := ctx.Get("db")
	if !ok {
		return nil
	}
	db, ok := idb.(*gorm.DB)
	if !ok {
		return nil
	}
	return db
}
