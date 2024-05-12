package UserController

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers"
	"github.com/iceking2nd/rustdesk-api-server/app/Middlewares/Database"
	"github.com/iceking2nd/rustdesk-api-server/app/Models"
	"github.com/iceking2nd/rustdesk-api-server/utils/Email"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"net/url"
)

// Resend godoc
//
//	@Summary	User activation token resend
//	@Schemes
//	@Description	User activation token resend
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			username	path	string	true	"Username which you registered"
//	@Success		204
//	@Failure		404	{object}	Controllers.ResponseError
//	@Failure		500	{object}	Controllers.ResponseError
//	@Router			/resend/{username} [get]
func Resend(c *gin.Context) {
	username := c.Param("username")
	db := Database.GetDB(c)
	var token Models.ActivateToken
	result := db.Preload(clause.Associations).First(&token, "username = ?", username)
	if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, Controllers.ResponseError{Error: "未找到相关的邮箱验证码"})
		c.Abort()
		return
	} else if result.Error != nil {
		c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: result.Error.Error()})
		c.Abort()
		return
	}
	subject := "账号激活邮件（重发）"
	content := `<p>%s，您好</p><br><p>您似乎丢了注册Rustdesk服务的账号时使用的邮箱验证码，请点击<a href="%s" target="_blank">此处</a>激活您的账号。如果您没有注册，请忽略此邮件。</p>`
	publicurl, err := url.Parse(viper.GetString("API.PublicURL"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: err.Error()})
		c.Abort()
		return
	}
	publicurl = publicurl.JoinPath("api", "activate", token.Token)
	mail := Email.NewMailProcessor()
	err = mail.Send(token.Username, subject, fmt.Sprintf(content, token.User.Name, publicurl.String()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
