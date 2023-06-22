package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.RouterGroup) {
	apiRoutesRegister(router)
	staticRoutesRegister(router)
}
