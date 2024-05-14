package UserController

import (
	"github.com/gin-gonic/gin"
	"github.com/iceking2nd/rustdesk-api-server/global"
	"net/http"
)

func LoginOptions(c *gin.Context) {
	log := global.Log.WithField("functions", "app.Controllers.UserController.LoginOptions")
	log.WithField("request", c.Request).Debugln("received request")
	c.JSON(http.StatusOK, []string{"common-oidc/[]"})
}
