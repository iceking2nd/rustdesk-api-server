package GroupController

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers"
	"github.com/iceking2nd/rustdesk-api-server/app/Middlewares/Database"
	"github.com/iceking2nd/rustdesk-api-server/app/Models"
	"github.com/tidwall/sjson"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"strings"
)

type GroupCreateRequest struct {
	AllowedIncomings []string `json:"allowed_incomings"`
	AllowedOutgoings []string `json:"allowed_outgoings"`
	Info             *struct {
		NoConnInGroup int `json:"no_conn_in_group"`
	} `json:"info,omitempty"`
	Name string  `json:"name"`
	Note *string `json:"note,omitempty"`
}

func Create(c *gin.Context) {
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

	var req GroupCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Controllers.ResponseError{Error: err.Error()})
		return
	}
	var group Models.Group
	group.GUID = uuid.New().String()
	group.Name = req.Name
	group.TeamID = token.User.Group.Team.ID
	if req.Note != nil {
		group.Note = *req.Note
	}
	group.Info = func() string {
		if req.Info != nil {
			info, err := sjson.Set("{}", "no_conn_in_group", req.Info.NoConnInGroup)
			if err != nil {
				c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: err.Error()})
				return "{}"
			}
			return info
		}
		return "{}"
	}()
	if len(req.AllowedIncomings) > 0 {
		for _, incoming := range req.AllowedIncomings {
			var grp Models.Group
			err = db.Where("name = ?", incoming).First(&grp).Error
			if err != nil {
				c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: err.Error()})
				return
			}
			group.AccessedFrom = append(group.AccessedFrom, &grp)
		}
	}
	if len(req.AllowedOutgoings) > 0 {
		for _, outgoing := range req.AllowedOutgoings {
			var grp Models.Group
			err = db.Where("name = ?", outgoing).First(&grp).Error
			if err != nil {
				c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: err.Error()})
				return
			}
			group.AccessTo = append(group.AccessTo, &grp)
		}
	}
	err = db.Create(&group).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}

type GroupUpdateRequest struct {
	AllowedIncomings *[]string `json:"allowed_incomings"`
	AllowedOutgoings *[]string `json:"allowed_outgoings"`
	Info             *struct {
		NoConnInGroup int `json:"no_conn_in_group"`
	} `json:"info,omitempty"`
	Name *string `json:"name,omitempty"`
	Note *string `json:"note,omitempty"`
	GUID string  `json:"guid"`
}

func Update(c *gin.Context) {
	var req GroupUpdateRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, Controllers.ResponseError{Error: err.Error()})
		return
	}
	db := Database.GetDB(c)
	var group Models.Group
	err = db.Preload(clause.Associations).Where("guid = ?", req.GUID).First(&group).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, Controllers.ResponseError{Error: err.Error()})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: err.Error()})
		return
	}

	switch {
	case req.Name != nil:
		group.Name = *req.Name
	case req.Info != nil:
		group.Info, err = sjson.Set(group.Info, "no_conn_in_group", req.Info.NoConnInGroup)
		if err != nil {
			c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: err.Error()})
			return
		}
	case req.Note != nil:
		group.Note = *req.Note
	case req.AllowedOutgoings != nil:
		group.AccessTo = nil
		if len(*req.AllowedOutgoings) > 0 {
			for _, outgoing := range *req.AllowedOutgoings {
				var grp Models.Group
				err = db.Where("name = ?", outgoing).First(&grp).Error
				if err != nil {
					c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: err.Error()})
					return
				}
				group.AccessTo = append(group.AccessTo, &grp)
			}
		}
	case req.AllowedIncomings != nil:
		group.AccessedFrom = nil
		if len(*req.AllowedIncomings) > 0 {
			for _, incoming := range *req.AllowedIncomings {
				var grp Models.Group
				err = db.Where("name = ?", incoming).First(&grp).Error
				if err != nil {
					c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: err.Error()})
					return
				}
				group.AccessedFrom = append(group.AccessedFrom, &grp)
			}
		}
	}
	err = db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&group).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: err.Error()})
		return
	}
	err = db.Session(&gorm.Session{FullSaveAssociations: true}).Model(&group).Association("AccessedFrom").Replace(group.AccessedFrom)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: err.Error()})
		return
	}
	err = db.Session(&gorm.Session{FullSaveAssociations: true}).Model(&group).Association("AccessTo").Replace(group.AccessTo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, nil)
}
