package AddressBookController

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/iceking2nd/rustdesk-api-server/app/Middlewares/Database"
	"github.com/iceking2nd/rustdesk-api-server/app/Models"
	"github.com/tidwall/sjson"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"strings"
)

func Get(c *gin.Context) {
	var (
		token Models.Token
		user  Models.User
		tags  []Models.Tag
		addrs []Models.Address
		data  = `{}`
	)
	db := Database.GetDB(c)
	err := db.Preload(clause.Associations).Where(&Models.Token{AccessToken: strings.Split(c.Request.Header.Get("Authorization"), " ")[1]}).First(&token).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Token已失效，请重新登录",
		})
		c.Abort()
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("查询Token数据时出现错误：%s", err.Error()),
		})
		c.Abort()
		return
	}
	user = token.User
	err = db.Where(&Models.Tag{UserID: user.ID}).Find(&tags).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("获取Tag数据时出现错误：%s", err.Error()),
		})
		c.Abort()
		return
	}
	if len(tags) > 0 {
		var t []string
		for _, tag := range tags {
			t = append(t, tag.Name)
		}
		data, err = sjson.Set(data, "tags", t)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("序列化Tag数据时出现错误：%s", err.Error()),
			})
			c.Abort()
			return
		}
	}
	err = db.Preload(clause.Associations).Where(&Models.Address{UserID: user.ID}).Find(&addrs).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("获取地址簿数据时出现错误：%s", err.Error()),
		})
		c.Abort()
		return
	}
	if len(addrs) > 0 {
		for _, addr := range addrs {
			a := make(map[string]any)
			a["id"] = addr.ClientID
			a["username"] = addr.Username
			a["hostname"] = addr.Hostname
			a["platform"] = addr.Platform
			a["alias"] = addr.Alias
			a["tags"] = make([]string, len(addr.Tags))
			for _, t := range addr.Tags {
				a["tags"] = append(a["tags"].([]string), t.Name)
			}
			data, err = sjson.Set(data, "peers.-1", a)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Sprintf("序列化地址簿数据时出现错误：%s", err.Error()),
				})
				c.Abort()
				return
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}
