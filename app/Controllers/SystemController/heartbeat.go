package SystemController

import (
	"github.com/gin-gonic/gin"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers"
	"net/http"
	"time"
)

type HeartBeatRequest struct {
	ID         string `json:"id"`
	ModifiedAt int64  `json:"modified_at"`
	UUID       string `json:"uuid"`
	Version    int64  `json:"ver"`
}

type HeartBeatResponse struct {
	ModifiedAt int64 `json:"modified_at"`
}

// HeartBeat godoc
// @Summary Heartbeat with clients
// @Schemes
// @Description Heartbeat with clients
// @Tags System
// @Accept json
// @Produce json
// @Success 200 {object} HeartBeatResponse "Always return unix timestamp of current time"
// @Router /heartbeat [post]
func HeartBeat(c *gin.Context) {
	var data HeartBeatRequest
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Controllers.ResponseError{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"modified_at": time.Now().Unix()})
}
