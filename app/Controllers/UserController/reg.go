package UserController

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers"
	"github.com/iceking2nd/rustdesk-api-server/app/Middlewares/Database"
	"github.com/iceking2nd/rustdesk-api-server/app/Models"
	"github.com/iceking2nd/rustdesk-api-server/utils/Email"
	"github.com/iceking2nd/rustdesk-api-server/utils/Hash"
	"github.com/spf13/viper"
	"net/http"
	"net/url"
)

type RegRequest struct {
	Username string `json:"username" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

// Reg Register godoc
//
//	@Summary	User registration
//	@Schemes
//	@Description	Register a new user
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			RegRequest	body	RegRequest	true	"Username must be email, password min length is 8, name is required"
//	@Success		204
//	@Failure		400	{object}	Controllers.ResponseError
//	@Failure		500	{object}	Controllers.ResponseError
//	@Router			/reg [post]
func Reg(c *gin.Context) {
	var data RegRequest
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, Controllers.ResponseError{Error: err.Error()})
		return
	}

	db := Database.GetDB(c)

	password, err := Hash.StringToSHA512(data.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: err.Error()})
		c.Abort()
		return
	}
	user := Models.User{Username: data.Username, Password: password, Name: data.Name, Status: 0}
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
	content := `<p>%s，您好</p><br><p>收到此邮件表示您正在注册Rustdesk服务的账号，请点击<a href="%s" target="_blank">此处</a>激活您的账号。如果您没有注册，请忽略此邮件。</p>`
	publicurl, err := url.Parse(viper.GetString("API.PublicURL"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: err.Error()})
		c.Abort()
		return
	}
	publicurl = publicurl.JoinPath("api", "activate", activateToken.Token)
	mail := Email.NewMailProcessor()
	err = mail.Send(user.Username, subject, fmt.Sprintf(content, user.Name, publicurl.String()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
