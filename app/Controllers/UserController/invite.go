package UserController

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iceking2nd/go-toolkits/generator"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers"
	"github.com/iceking2nd/rustdesk-api-server/app/Middlewares/Database"
	"github.com/iceking2nd/rustdesk-api-server/app/Models"
	"github.com/iceking2nd/rustdesk-api-server/utils/Email"
	"github.com/iceking2nd/rustdesk-api-server/utils/Hash"
	"github.com/spf13/viper"
	"net/http"
	"net/url"
)

type InviteRequest struct {
	EMail     string `json:"email" binding:"required,email"`
	GroupName string `json:"group_name" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Note      string `json:"note,omitempty"`
}

// Invite godoc
//
//	@Summary	User registration
//	@Schemes
//	@Description	Register a new user
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			InviteRequest	body	InviteRequest	true
//	@Success		204
//	@Failure		400	{object}	Controllers.ResponseError
//	@Failure		500	{object}	Controllers.ResponseError
//	@Router			/invite [post]
func Invite(c *gin.Context) {
	genPassword := generator.RandStringBytesMaskImprSrcUnsafe(30)
	var data InviteRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, Controllers.ResponseError{Error: err.Error()})
		return
	}

	db := Database.GetDB(c)

	password, err := Hash.StringToSHA512(genPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: err.Error()})
		c.Abort()
		return
	}
	user := Models.User{Username: data.EMail, Password: password, Name: data.Name, Status: 0, Note: data.Note, IsAdmin: false, GUID: uuid.New().String()}
	result := db.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, Controllers.ResponseError{Error: result.Error.Error()})
		c.Abort()
		return
	}
	activateToken := Models.ActivateToken{
		Token:    uuid.New().String(),
		Username: user.Username,
		UserID:   user.ID,
	}
	result = db.Create(&activateToken)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, Controllers.ResponseError{Error: result.Error.Error()})
		c.Abort()
		return
	}
	subject := "账号激活邮件"
	content := `<p>%s，您好</p><br>
<p>收到此邮件表示您正在注册Rustdesk服务的账号，请点击<a href="%s" target="_blank">此处</a>激活您的账号。如果您没有注册，请忽略此邮件。</p><br>

<p>登录信息</p><br>
<p>登录用户名：%s</p><br>
<p>登录密码：%s</p>`
	publicurl, err := url.Parse(viper.GetString("API.PublicURL"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: err.Error()})
		c.Abort()
		return
	}
	publicurl = publicurl.JoinPath("api", "activate", activateToken.Token)
	mail := Email.NewMailProcessor()
	err = mail.Send(user.Username, subject, fmt.Sprintf(content, user.Name, publicurl.String(), user.Username, genPassword))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func ReInvite(c *gin.Context) {}
