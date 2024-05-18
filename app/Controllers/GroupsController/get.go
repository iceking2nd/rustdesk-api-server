package GroupsController

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/iceking2nd/rustdesk-api-server/app/Middlewares/Database"
	"github.com/iceking2nd/rustdesk-api-server/app/Models"
	"github.com/iceking2nd/rustdesk-api-server/global"
	"github.com/iceking2nd/rustdesk-api-server/utils/Paginate"
	"net/http"
	"strconv"
)

type GroupsListResponse struct {
	Data  []map[string]interface{} `json:"data"`
	Total int64                    `json:"total"`
}

func List(c *gin.Context) {
	db := Database.GetDB(c)
	page, _ := strconv.Atoi(c.DefaultQuery("current", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	orderby := c.DefaultQuery("orderby", "id")
	sort := c.DefaultQuery("sort", "desc")
	pg := Paginate.Paging(&Paginate.PagingParam{
		DB:      db,
		Page:    page,
		Limit:   limit,
		OrderBy: []string{fmt.Sprintf("%s %s", orderby, sort)},
		ShowSQL: func() bool {
			if global.LogLevel >= 5 {
				return true
			}
			return false
		}(),
	}, &[]Models.Group{})
	c.JSON(http.StatusOK, groupsListPageinateWrapper(pg, c))
}

func groupsListPageinateWrapper(model interface{}, c *gin.Context) interface{} {
	ul := new(GroupsListResponse)
	var data []map[string]interface{}
	for _, d := range *model.(*Paginate.Paginator).Data.(*[]Models.Group) {
		data = append(data, map[string]interface{}{
			"guid":       d.GUID,
			"name":       d.Name,
			"note":       d.Note,
			"team":       d.Team.GUID,
			"created_at": d.CreatedAt.Unix(),
			"access_to": func(accessTo []*Models.Group) []interface{} {
				if len(accessTo) == 0 {
					return []interface{}{}
				}
				var grps []interface{}
				for _, v := range accessTo {
					grps = append(grps, v)
				}
				return grps
			}(d.AccessTo),
			"accessed_from": func(accessedFrom []*Models.Group) []interface{} {
				if len(accessedFrom) == 0 {
					return []interface{}{}
				}
				var grps []interface{}
				for _, v := range accessedFrom {
					grps = append(grps, v)
				}
				return grps
			}(d.AccessedFrom),
		})
	}
	ul.Data = data
	ul.Total = model.(*Paginate.Paginator).TotalRecord
	return ul
}
