package UserController

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers"
	"github.com/iceking2nd/rustdesk-api-server/app/Middlewares/Database"
	"github.com/iceking2nd/rustdesk-api-server/app/Models"
	"github.com/iceking2nd/rustdesk-api-server/utils/Hash"
	"github.com/tidwall/sjson"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"strings"
)

func Post(c *gin.Context) {}

type UserPutRequest struct {
	EMail           *string `form:"email" json:"email,omitempty" binding:"omitempty,email"`
	Password        *string `json:"password,omitempty"`
	ConfirmPassword *string `json:"confirm_password,omitempty"`
	Info            *struct {
		EMailVerification      *bool `json:"email_verification,omitempty"`
		EMailAlarmNotification *bool `json:"email_alarm_notification,omitempty"`
	} `json:"info,omitempty"`
}

func Put(c *gin.Context) {
	rt := strings.Split(c.Request.Header.Get("Authorization"), " ")[1]
	db := Database.GetDB(c)
	var token Models.Token
	err := db.Preload(clause.Associations).Where(&Models.Token{AccessToken: rt}).First(&token).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusUnauthorized, Controllers.ResponseError{Error: "Invalid token"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: fmt.Sprintf("查询Token数据时出现错误：%s", err.Error())})
		return
	}

	user := token.User
	var reqData UserPutRequest
	err = c.ShouldBindJSON(&reqData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: err.Error()})
		return
	}
	if reqData.EMail != nil {
		user.Username = *reqData.EMail
		err = db.Save(&user).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: err.Error()})
			return
		}
	}
	if reqData.Password != nil && reqData.ConfirmPassword != nil && *reqData.ConfirmPassword == *reqData.Password {
		hashedPassword, err := Hash.StringToSHA512(*reqData.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: err.Error()})
			return
		}
		user.Password = hashedPassword
		err = db.Save(&user).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: err.Error()})
			return
		}
	}
	if reqData.Info != nil {
		userInfo := user.Info
		if reqData.Info.EMailVerification != nil {
			userInfo, err = sjson.Set(userInfo, "email_verification", *reqData.Info.EMailVerification)
			if err != nil {
				c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: err.Error()})
				return
			}
		}
		if reqData.Info.EMailAlarmNotification != nil {
			userInfo, err = sjson.Set(userInfo, "email_alarm_notification", *reqData.Info.EMailAlarmNotification)
			if err != nil {
				c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: err.Error()})
				return
			}
		}
		user.Info = userInfo
		err = db.Save(&user).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: err.Error()})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"verification_email_sent": false,
	})
}
func Delete(c *gin.Context)           {}
func DeleteSession(c *gin.Context)    {}
func UnsubscribeEmail(c *gin.Context) {}
func UpdateVerify(c *gin.Context)     {}
func VerificationCode(c *gin.Context) {}
