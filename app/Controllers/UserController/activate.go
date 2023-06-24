package UserController

import (
	"github.com/gin-gonic/gin"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers"
	"github.com/iceking2nd/rustdesk-api-server/app/Middlewares/Database"
	"github.com/iceking2nd/rustdesk-api-server/app/Models"
	"gorm.io/gorm/clause"
	"net/http"
)

// Activate godoc
// @Summary User account activation
// @Schemes
// @Description User account activation
// @Tags User
// @Accept json
// @Produce json
// @Param token path string true "Activation token user received by email"
// @Success 204
// @Failure 500 {object} Controllers.ResponseError
// @Router /activate/{token} [get]
func Activate(c *gin.Context) {
	token := c.Param("token")
	db := Database.GetDB(c)
	var activateToken Models.ActivateToken
	result := db.Preload(clause.Associations).First(&activateToken, "token = ?", token)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: result.Error.Error()})
		c.Abort()
		return
	}
	user := activateToken.User
	user.IsValidated = true
	result = db.Updates(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: result.Error.Error()})
		c.Abort()
		return
	}
	result = db.Delete(&activateToken)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: result.Error.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
