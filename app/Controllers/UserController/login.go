package UserController

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers"
	"github.com/iceking2nd/rustdesk-api-server/app/Middlewares/Database"
	"github.com/iceking2nd/rustdesk-api-server/app/Models"
	"github.com/iceking2nd/rustdesk-api-server/global"
	"github.com/iceking2nd/rustdesk-api-server/utils/Hash"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type LoginRequest struct {
	Username   string                 `json:"username"`
	Password   string                 `json:"password"`
	ClientID   string                 `json:"id"`
	UUID       string                 `json:"uuid"`
	AutoLogin  bool                   `json:"autoLogin"`
	Type       string                 `json:"type"`
	DeviceInfo LoginRequestDeviceInfo `json:"deviceInfo"`
}

type LoginRequestDeviceInfo struct {
	OS   string `json:"os"`
	Type string `json:"type"`
	Name string `json:"name"`
}

type LoginResponse struct {
	Type        string            `json:"type"`
	AccessToken string            `json:"access_token"`
	TFAType     string            `json:"tfa_type"`
	Secret      string            `json:"secret"`
	User        LoginResponseUser `json:"user"`
}

type LoginResponseUser struct {
	Name    string                `json:"name"`
	Note    string                `json:"note"`
	Info    LoginResponseUserInfo `json:"info"`
	IsAdmin bool                  `json:"is_admin"`
}

type LoginResponseUserInfo struct {
	EmailAlarmNotification bool `json:"email_alarm_notification"`
}

// Login godoc
//
//	@Summary	User login
//	@Schemes
//	@Description	User login
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			LoginRequest	body		LoginRequest	true	"Login data"
//	@Success		200				{object}	LoginResponse
//	@Failure		403				{object}	Controllers.ResponseError
//	@Failure		500				{object}	Controllers.ResponseError
//	@Router			/login [post]
func Login(c *gin.Context) {
	log := global.Log.WithField("functions", "app.Controllers.UserController.Login")
	db := Database.GetDB(c)
	var data LoginRequest
	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: fmt.Sprintf("接收登录数据时出现错误：%s", err.Error())})
		return
	}
	log.WithField("data", data).Debugln("Login Request Data")
	hashedPassword, err := Hash.StringToSHA512(data.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: fmt.Sprintf("处理密码时出现错误：%s", err.Error())})
		return
	}
	var u Models.User
	err = db.Where(&Models.User{Username: data.Username, Password: hashedPassword, IsValidated: true}).First(&u).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusForbidden, Controllers.ResponseError{Error: "您的帐号、密码不正确或帐号未经过验证。"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: fmt.Sprintf("登录验证时出现错误：%s", err.Error())})
		return
	}
	token, err := Hash.StringToMD5(fmt.Sprintf("%s|%s|%s|%d", data.Username, data.ClientID, data.UUID, time.Now().UnixNano()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: fmt.Sprintf("生成Token时出现错误：%s", err.Error())})
		return
	}
	t := Models.Token{AccessToken: token, UserID: u.ID, ClientID: data.ClientID}
	result := db.Create(&t)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: fmt.Sprintf("保存Token时出现错误：%s", result.Error.Error())})
		return
	}

	var addr Models.Address
	err = db.Where(&Models.Address{UserID: u.ID, ClientID: data.ClientID}).FirstOrCreate(&addr).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: fmt.Sprintf("向地址簿保存本机信息时出现错误：%s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, LoginResponse{
		Type:        "access_token",
		AccessToken: token,
		TFAType:     "",
		Secret:      "",
		User: LoginResponseUser{
			Name: u.Username,
			Note: "",
			Info: LoginResponseUserInfo{
				EmailAlarmNotification: true,
			},
			IsAdmin: true,
		},
	})
}
