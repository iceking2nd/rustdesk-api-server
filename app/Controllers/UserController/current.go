package UserController

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers"
	"github.com/iceking2nd/rustdesk-api-server/app/Middlewares/Database"
	"github.com/iceking2nd/rustdesk-api-server/app/Models"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"strings"
)

type CurrentUserResponse struct {
	Name    string      `json:"name"`
	IsAdmin bool        `json:"is_admin"`
	Note    string      `json:"note"`
	Info    interface{} `json:"info"`
	EMail   string      `json:"email"`
}

// CurrentUser godoc
//
//	@Summary	CurrentUser Info
//	@Schemes
//	@Description	Get current login user info
//	@Tags			User
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	CurrentUserResponse
//	@Failure		404	{object}	Controllers.ResponseError
//	@Failure		500	{object}	Controllers.ResponseError
//	@Router			/currentUser [post]
func CurrentUser(c *gin.Context) {
	rt := strings.Split(c.Request.Header.Get("Authorization"), " ")[1]
	db := Database.GetDB(c)
	var token Models.Token
	err := db.Preload(clause.Associations).Where(&Models.Token{AccessToken: rt}).First(&token).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusUnauthorized, Controllers.ResponseError{Error: "Invalid token"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: fmt.Sprintf("查询Token数据时出现错误：%s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, CurrentUserResponse{
		Name:    token.User.Name,
		IsAdmin: token.User.IsAdmin,
		Note:    token.User.Note,
		Info: func(info string) map[string]interface{} {
			i := make(map[string]interface{})
			gjson.Parse(info).ForEach(func(key, value gjson.Result) bool {
				i[key.String()] = value.Value()
				return true
			})
			return i
		}(token.User.Info),
		EMail: token.User.Username,
	})
}
