package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers/AddressBookController"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers/SystemController"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers/UserController"
	"github.com/iceking2nd/rustdesk-api-server/app/Middlewares/Auth"
)

func apiRoutesRegister(route *gin.RouterGroup) {
	apiRoutes := route.Group("/api")
	apiRoutes.POST("/reg", UserController.Reg)
	apiRoutes.GET("/activate/:token", UserController.Activate)
	apiRoutes.GET("/resend/:username", UserController.Resend)
	apiRoutes.POST("/login", UserController.Login)
	apiRoutes.POST("/logout", UserController.Logout).Use(Auth.TokenAuth())
	apiRoutes.POST("/currentUser", UserController.CurrentUser).Use(Auth.TokenAuth())
	apiRoutes.GET("/audit", SystemController.Audit).Use(Auth.TokenAuth())
	apiRoutes.POST("/audit", SystemController.Audit).Use(Auth.TokenAuth())
	apiRoutes.GET("/heartbeat", SystemController.HeartBeat).Use(Auth.TokenAuth())
	apiRoutes.POST("/heartbeat", SystemController.HeartBeat).Use(Auth.TokenAuth())

	AddressBookRoutes := apiRoutes.Group("/ab").Use(Auth.TokenAuth())
	AddressBookRoutes.POST("", AddressBookController.Update)
	AddressBookRoutes.POST("/get", AddressBookController.Get)
}
