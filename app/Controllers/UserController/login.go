package UserController

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/iceking2nd/rustdesk-api-server/app/Middlewares/Database"
	"github.com/iceking2nd/rustdesk-api-server/app/Models"
	"github.com/iceking2nd/rustdesk-api-server/utils/Hash"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	ClientID string `json:"id"`
	UUID     string `json:"uuid"`
}

func Login(c *gin.Context) {
	db := Database.GetDB(c)
	var data LoginRequest
	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("接收登录数据时出现错误：%s", err.Error()),
		})
		return
	}
	hashedPassword, err := Hash.StringToSHA512(data.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("处理密码时出现错误：%s", err.Error()),
		})
		return
	}
	var u Models.User
	err = db.Where(&Models.User{Username: data.Username, Password: hashedPassword, IsValidated: true}).First(&u).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "您的帐号、密码不正确或帐号未经过验证。",
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("登录验证时出现错误：%s", err.Error()),
		})
		return
	}
	token, err := Hash.StringToMD5(fmt.Sprintf("%s|%s|%s|%d", data.Username, data.ClientID, data.UUID, time.Now().UnixNano()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("生成Token时出现错误：%s", err.Error()),
		})
		return
	}
	t := Models.Token{AccessToken: token, UserID: u.ID, ClientID: data.ClientID}
	result := db.Create(&t)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("保存Token时出现错误：%s", result.Error.Error()),
		})
		return
	}

	var addr Models.Address
	err = db.Where(&Models.Address{UserID: u.ID, ClientID: data.ClientID}).FirstOrCreate(&addr).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("向地址簿保存本机信息时出现错误：%s", err.Error()),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"type":         "access_token",
		"access_token": token,
		"user": gin.H{
			"name": u.Username,
		},
	})
}
