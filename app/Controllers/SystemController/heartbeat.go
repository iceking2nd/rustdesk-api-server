package SystemController

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func HeartBeat(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"modified_at": time.Now().Unix()})
}
