package SystemController

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type HeartBeatResponse struct {
	ModifiedAt int64 `json:"modified_at"`
}

// HeartBeat godoc
// @Summary Heartbeat with clients
// @Schemes
// @Description Heartbeat with clients
// @Tags System
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} HeartBeatResponse "Always return unix timestamp of current time"
// @Router /heartbeat [get]
// @Router /heartbeat [post]
func HeartBeat(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"modified_at": time.Now().Unix()})
}
