package UserController

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/iceking2nd/rustdesk-api-server/app/Middlewares/Database"
	"github.com/iceking2nd/rustdesk-api-server/app/Models"
	"github.com/iceking2nd/rustdesk-api-server/utils/Paginate"
	"net/http"
	"strconv"
)

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
		ShowSQL: true,
	}, &Models.User{})
	c.JSON(http.StatusOK, Paginate.PageinateWrapper(pg, c))
}
