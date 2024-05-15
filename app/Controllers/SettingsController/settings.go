package SettingsController

import (
	"github.com/gin-gonic/gin"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers"
	"github.com/iceking2nd/rustdesk-api-server/app/Middlewares/Database"
	"github.com/iceking2nd/rustdesk-api-server/app/Models"
	"net/http"
)

func Put(c *gin.Context) {
	var data map[string]string
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: err.Error()})
		return
	}
	db := Database.GetDB(c)

	for k, v := range data {
		setting := Models.Settings{Key: k, Value: v}
		db.Where(Models.Settings{Key: k}).Assign(Models.Settings{Value: v}).FirstOrCreate(&setting)
	}
	c.Writer.WriteHeader(http.StatusOK)
}
