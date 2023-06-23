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

func Logout(c *gin.Context) {
	db := Database.GetDB(c).Session(&gorm.Session{FullSaveAssociations: true})
	var token Models.Token
	err := db.Preload(clause.Associations).Where(&Models.Token{AccessToken: strings.Split(c.Request.Header.Get("Authorization"), " ")[1]}).First(&token).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Token已失效，无需登出",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("查询Token数据时出现错误：%s", err.Error()),
		})
		return
	}

	var addr Models.Address

	err = db.Where("client_id = ? AND user_id = ?", token.ClientID, token.UserID).First(&addr).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("获取终端对应地址簿条目时出现错误：%s", err.Error()),
		})
		return
	}
	err = db.Select(clause.Associations).Delete(&addr).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("清理终端对应地址簿条目时出现错误：%s", err.Error()),
		})
		return
	}
	err = db.Delete(&token).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("清理Token时出现错误：%s", err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "登出成功"})
}
