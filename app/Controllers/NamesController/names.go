package NamesController

import (
	"github.com/gin-gonic/gin"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers"
	"github.com/iceking2nd/rustdesk-api-server/app/Middlewares/Database"
	"github.com/iceking2nd/rustdesk-api-server/app/Models"
	"net/http"
)

func Get(c *gin.Context) {
	db := Database.GetDB(c)
	table := c.DefaultQuery("table", "grp")
	switch table {
	case "grp":
		var groups []Models.Group
		err := db.Find(&groups).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: err.Error()})
			return
		}
		var groupsName []string
		for _, group := range groups {
			if len(group.Name) > 0 {
				groupsName = append(groupsName, group.Name)
			}
		}
		c.JSON(http.StatusOK, groupsName)
	}
}
