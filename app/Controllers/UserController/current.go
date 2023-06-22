package UserController

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/iceking2nd/rustdesk-api-server/app/Middlewares/Database"
	"github.com/iceking2nd/rustdesk-api-server/app/Models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"strings"
)

func CurrentUser(c *gin.Context) {
	rt := strings.Split(c.Request.Header.Get("Authorization"), " ")[1]
	db := Database.GetDB(c)
	var token Models.Token
	err := db.Preload(clause.Associations).Where(&Models.Token{AccessToken: rt}).First(&token).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Token已失效，请重新登录",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("查询Token数据时出现错误：%s", err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"name": token.User.Username,
	})
}
