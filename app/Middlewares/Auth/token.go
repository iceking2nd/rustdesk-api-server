package Auth

import (
	"github.com/gin-gonic/gin"
	"github.com/iceking2nd/rustdesk-api-server/app/Middlewares/Database"
	"github.com/iceking2nd/rustdesk-api-server/app/Models"
	"github.com/iceking2nd/rustdesk-api-server/global"
	"net/http"
	"strings"
)

func TokenAuth() gin.HandlerFunc {
	log := global.Log.WithField("function", "app.Middlewares.Auth.TokenAuth")
	return func(c *gin.Context) {
		log.WithField("request", c.Request).Debugln("received request")
		auth := strings.Split(c.Request.Header.Get("Authorization"), " ")
		if len(auth) != 2 || auth[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证错误，认证方式不正确。"})
			c.Abort()
			return
		}
		db := Database.GetDB(c)
		var token Models.Token
		err := db.Where(&Models.Token{AccessToken: auth[1]}).First(&token).Error
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证错误，Token不存在。"})
			c.Abort()
			return
		}
		c.Next()
	}
}
