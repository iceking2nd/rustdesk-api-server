package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/iceking2nd/rustdesk-api-server/frontend"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"path/filepath"
)

func staticRoutesRegister(router *gin.RouterGroup) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.StaticFS("/icons", frontend.GetFS(frontend.IconsFS, "icons"))
	router.StaticFS("/scripts", frontend.GetFS(frontend.ScriptsFS, "scripts"))
	router.StaticFS("/styles", frontend.GetFS(frontend.StylesFS, "styles"))
	router.GET("/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		switch filepath.Ext(filename) {
		case ".svg":
			c.Writer.Header().Set("Content-Type", "image/svg+xml")
		case ".png":
			c.Writer.Header().Set("Content-Type", "image/png")
		case ".js":
			c.Writer.Header().Set("Content-Type", "application/javascript")
		case ".css":
			c.Writer.Header().Set("Content-Type", "text/css")
		default:
			c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		}
		c.Writer.WriteHeader(200)
		index, _ := frontend.StaticFS.ReadFile(fmt.Sprintf("static/%s", filename))
		_, _ = c.Writer.Write(index)
	})
}
