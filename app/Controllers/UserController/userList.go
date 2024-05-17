package UserController

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

type UserListResponse struct {
	Data  []map[string]interface{} `json:"data"`
	Total int64                    `json:"total"`
}

func UserList(c *gin.Context) {
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
	}, &[]Models.User{})
	c.JSON(http.StatusOK, userListPageinateWrapper(pg, c))
}

func userListPageinateWrapper(model interface{}, c *gin.Context) interface{} {
	ul := new(UserListResponse)
	var data []map[string]interface{}
	for _, d := range *model.(*Paginate.Paginator).Data.(*[]Models.User) {
		data = append(data, map[string]interface{}{
			"guid":       d.GUID,
			"name":       d.Name,
			"note":       d.Note,
			"email":      d.Username,
			"status":     d.Status,
			"group_name": "Default",
			"is_admin":   d.IsAdmin,
		})
	}
	ul.Data = data
	ul.Total = model.(*Paginate.Paginator).TotalRecord
	return ul
}
