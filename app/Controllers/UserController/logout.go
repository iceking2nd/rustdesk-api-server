package UserController

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers"
	"github.com/iceking2nd/rustdesk-api-server/app/Middlewares/Database"
	"github.com/iceking2nd/rustdesk-api-server/app/Models"
	"github.com/iceking2nd/rustdesk-api-server/global"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"strings"
)

// Logout godoc
//
//	@Summary	User logout
//	@Schemes
//	@Description	User logout
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	Controllers.Response
//	@Failure		404	{object}	Controllers.ResponseError
//	@Failure		500	{object}	Controllers.ResponseError
//	@Router			/logout [post]
func Logout(c *gin.Context) {
	log := global.Log.WithField("functions", "app.Controllers.UserController.Logout")
	log.WithField("request", c.Request).Debugln("Logout request")
	db := Database.GetDB(c).Session(&gorm.Session{FullSaveAssociations: true})
	var token Models.Token
	err := db.Preload(clause.Associations).Where(&Models.Token{AccessToken: strings.Split(c.Request.Header.Get("Authorization"), " ")[1]}).First(&token).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, Controllers.ResponseError{Error: "Token已失效，无需登出"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: fmt.Sprintf("查询Token数据时出现错误：%s", err.Error())})
		return
	}

	if token.LoginTokenType != "browser" {
		var addr Models.Address

		err = db.Where("client_id = ? AND user_id = ?", token.ClientID, token.UserID).First(&addr).Error
		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: fmt.Sprintf("获取终端对应地址簿条目时出现错误：%s", err.Error())})
			return
		}
		err = db.Select(clause.Associations).Delete(&addr).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: fmt.Sprintf("清理终端对应地址簿条目时出现错误：%s", err.Error())})
			return
		}
	}
	err = db.Delete(&token).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: fmt.Sprintf("清理Token时出现错误：%s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, Controllers.Response{Data: "登出成功"})
}
