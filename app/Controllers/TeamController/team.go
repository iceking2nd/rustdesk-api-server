package TeamController

import (
	"github.com/gin-gonic/gin"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers"
	"github.com/spf13/viper"
	"net/http"
)

func InfoGet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

type TeamInfoUpdateRequest struct {
	SMTP *struct {
		Host         string `form:"host" json:"host,omitempty" binding:"required,omitempty"`
		Port         int    `form:"port" json:"port,omitempty" binding:"omitempty,numeric,min=1,max=65535"`
		From         string `form:"from" json:"from,omitempty" binding:"required,omitempty"`
		RequiresAuth bool   `form:"require_auth" json:"requires_auth,omitempty" binding:"boolean,omitempty"`
		Account      string `form:"account" json:"account,omitempty" binding:"required,omitempty"`
		Password     string `form:"password" json:"password,omitempty" binding:"required,omitempty"`
		PermFile     string `form:"permfile" json:"permfile,omitempty"`
		StartTLS     bool   `form:"starttls" json:"starttls,omitempty" binding:"boolean,omitempty"`
	} `json:"smtp,omitempty"`
}

func InfoPut(c *gin.Context) {
	var reqData TeamInfoUpdateRequest
	err := c.ShouldBindJSON(&reqData)
	if err != nil {
		c.JSON(http.StatusOK, Controllers.ResponseError{Error: err.Error()})
		return
	}
	if reqData.SMTP != nil {
		viper.Set("SMTP.From", reqData.SMTP.From)
		viper.Set("SMTP.Host", reqData.SMTP.Host)
		viper.Set("SMTP.Port", reqData.SMTP.Port)
		viper.Set("SMTP.Username", reqData.SMTP.Account)
		viper.Set("SMTP.Password", reqData.SMTP.Password)
		viper.Set("SMTP.RequiresAuth", reqData.SMTP.RequiresAuth)
		viper.Set("SMTP.PermFile", reqData.SMTP.PermFile)
		viper.Set("SMTP.StartTLS", reqData.SMTP.StartTLS)
		err = viper.WriteConfig()
		if err != nil {
			c.JSON(http.StatusOK, Controllers.ResponseError{Error: err.Error()})
		}
	}
}
