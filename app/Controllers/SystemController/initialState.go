package SystemController

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers"
	"github.com/iceking2nd/rustdesk-api-server/app/Middlewares/Database"
	"github.com/iceking2nd/rustdesk-api-server/app/Models"
	"github.com/iceking2nd/rustdesk-api-server/global"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
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
			EMailVerification      bool `json:"email_verification"`
			EmailAlarmNotification bool `json:"email_alarm_notification"`
		} `json:"info"`
		IsAdmin bool   `json:"is_admin"`
		EMail   string `json:"email"`
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
	r.User.Note = token.User.Note
	r.User.EMail = token.User.Username
	r.User.IsAdmin = token.User.IsAdmin
	userInfo := gjson.Parse(token.User.Info)
	r.User.Info.EMailVerification = func() bool {
		if userInfo.Get("email_verification").Exists() {
			return userInfo.Get("email_verification").Bool()
		}
		return false
	}()
	r.User.Info.EmailAlarmNotification = func() bool {
		if userInfo.Get("email_alarm_notification").Exists() {
			return userInfo.Get("email_alarm_notification").Bool()
		}
		return false
	}()

	var settings []Models.Settings
	err := db.Find(&settings).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: err.Error()})
		return
	}
	for _, setting := range settings {
		r.Settings = append(r.Settings, struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}{Key: setting.Key, Value: setting.Value})
	}
	c.JSON(http.StatusOK, r)
}
