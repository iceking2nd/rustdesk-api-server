package UserController

import (
	"github.com/gin-gonic/gin"
	"github.com/iceking2nd/rustdesk-api-server/app/Middlewares/Database"
	"github.com/iceking2nd/rustdesk-api-server/app/Models"
	"gorm.io/gorm/clause"
	"net/http"
)

func Activate(c *gin.Context) {
	token := c.Param("token")
	db := Database.GetDB(c)
	var activateToken Models.ActivateToken
	result := db.Preload(clause.Associations).First(&activateToken, "token = ?", token)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		c.Abort()
		return
	}
	user := activateToken.User
	user.IsValidated = true
	result = db.Updates(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		c.Abort()
		return
	}
	result = db.Delete(&activateToken)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
