package SMTPController

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
)

func Validate(c *gin.Context) {
	c.JSON(http.StatusOK, viper.IsSet("SMTP.Host") && viper.IsSet("SMTP.Port") && viper.IsSet("SMTP.Username") && viper.IsSet("SMTP.Password"))
}
