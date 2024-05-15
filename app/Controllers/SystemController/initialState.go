package SystemController

import (
	"github.com/gin-gonic/gin"
	"github.com/iceking2nd/rustdesk-api-server/app/Middlewares/Database"
	"github.com/iceking2nd/rustdesk-api-server/app/Models"
	"github.com/iceking2nd/rustdesk-api-server/global"
	"gorm.io/gorm/clause"
	"net/http"
	"strings"
)

type InitialStateResponse struct {
	Version string `json:"version"`
	User    struct {
		Name string `json:"name"`
		Note string `json:"note"`
		Info struct {
			EmailAlarmNotification bool `json:"email_alarm_notification"`
		} `json:"info"`
		IsAdmin bool `json:"is_admin"`
	} `json:"user"`
	Team struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"team"`
	Group struct {
		Guid      string `json:"guid"`
		Name      string `json:"name"`
		Team      string `json:"team"`
		CreatedAt int    `json:"created_at"`
	} `json:"group"`
	License struct {
		StatusText string `json:"status_text"`
	} `json:"license"`
	Settings []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"settings"`
}

func InitialState(c *gin.Context) {
	db := Database.GetDB(c)
	r := new(InitialStateResponse)
	if len(global.Version) > 0 {
		r.Version = "rustdesk-api-server " + global.Version
	} else {
		r.Version = "rustdesk-api-server " + "undefined"
	}
	r.License.StatusText = ""

	var token Models.Token
	db.Preload(clause.Associations).First(&token, "access_token = ?", strings.Split(c.Request.Header.Get("Authorization"), " ")[1])
	r.User.Name = token.User.Name
	r.User.IsAdmin = token.User.IsAdmin

	var settings []Models.Settings
	db.Find(&settings)
	for _, setting := range settings {
		r.Settings = append(r.Settings, struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}{Key: setting.Key, Value: setting.Value})
	}
	c.JSON(http.StatusOK, r)
}
