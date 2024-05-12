package AddressBookController

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

type AddressBookUpdateRequest struct {
	Data string `json:"data"`
}

// Update address book and tags godoc
//
//	@Summary	Update all address book and tags data
//	@Schemes
//	@Description	Update all address book and tags data
//	@Tags			Address book and tags
//	@Security		ApiKeyAuth
//	@Accept			json
//	@Produce		json
//	@Param			AddressBookUpdateRequest	body	AddressBookUpdateRequest	true	"Update data is a serialized json string"
//	@Success		200
//	@Failure		default	{object}	Controllers.ResponseError
//	@Router			/ab [post]
func Update(c *gin.Context) {
	var (
		reqdata AddressBookUpdateRequest
		user    Models.User
		token   Models.Token
	)
	db := Database.GetDB(c).Session(&gorm.Session{FullSaveAssociations: true})

	err := db.Preload(clause.Associations).Where(&Models.Token{AccessToken: strings.Split(c.Request.Header.Get("Authorization"), " ")[1]}).First(&token).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, Controllers.ResponseError{Error: "Token已失效，请重新登录"})
		c.Abort()
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: fmt.Sprintf("查询Token数据时出现错误：%s", err.Error())})
		c.Abort()
		return
	}
	user = token.User

	err = c.BindJSON(&reqdata)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: fmt.Sprintf("接收地址簿数据时出现错误：%s", err.Error())})
		c.Abort()
		return
	}
	data := gjson.Parse(reqdata.Data)
	data.Get("tags").ForEach(func(key, value gjson.Result) bool {
		var tag Models.Tag
		err = db.Where(Models.Tag{Name: value.String()}).Assign(Models.Tag{UserID: user.ID}).FirstOrCreate(&tag).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: fmt.Sprintf("保存Tag数据时出现错误：%s", err.Error())})
			c.Abort()
			return false
		}
		return true
	})
	data.Get("peers").ForEach(func(key, value gjson.Result) bool {
		var (
			tags []Models.Tag
			addr Models.Address
			t    []string
		)
		value.Get("tags").ForEach(func(key, value gjson.Result) bool {
			t = append(t, value.String())
			return true
		})
		if len(t) > 0 {
			result := db.Where("name IN ?", t).Find(&tags)
			if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: fmt.Sprintf("查询Tag数据时出现错误：%s", err.Error())})
				c.Abort()
				return false
			}
		}

		address := Models.Address{
			UserID:   user.ID,
			ClientID: value.Get("id").String(),
			Hostname: value.Get("hostname").String(),
			Username: value.Get("username").String(),
			Platform: value.Get("platform").String(),
			Alias:    value.Get("alias").String(),
			Tags:     tags,
		}
		result := db.Preload(clause.Associations).Where(Models.Address{UserID: user.ID, ClientID: value.Get("id").String()})
		result = result.Assign(address).FirstOrCreate(&addr)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: fmt.Sprintf("保存地址簿数据时出现错误：%s", result.Error.Error())})
			c.Abort()
			return false
		}
		addr.Tags = tags
		err = db.Model(&addr).Association("Tags").Replace(tags)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: fmt.Sprintf("保存地址Tag时出现错误：%s", result.Error.Error())})
			c.Abort()
			return false
		}
		return true
	})
	c.JSON(http.StatusOK, nil)
}
