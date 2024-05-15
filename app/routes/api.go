package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers/AddressBookController"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers/AuditController"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers/AuditsController"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers/CrossGroupController"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers/CustomClientsController"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers/DevicesController"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers/GeoController"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers/GroupController"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers/GroupsController"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers/KeypairController"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers/NamesController"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers/OIDCController"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers/PeerController"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers/PeersController"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers/SMTPController"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers/SessionController"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers/SettingsController"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers/StrategiesController"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers/StrategyController"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers/SystemController"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers/TeamController"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers/TokensController"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers/UserController"
	"github.com/iceking2nd/rustdesk-api-server/app/Controllers/VerifyController"
	"github.com/iceking2nd/rustdesk-api-server/app/Middlewares/Auth"
)

func apiRoutesRegister(route *gin.RouterGroup) {
	apiRoutes := route.Group("/api")
	apiRoutes.POST("/login", UserController.Login)
	apiRoutes.GET("/login-options", UserController.LoginOptions)
	apiRoutes.GET("/activate/:token", UserController.Activate)
	apiRoutes.GET("/cfg", SystemController.Cfg)
	apiRoutes.GET("/check/read", SystemController.CheckRead)
	apiRoutes.PUT("/delete-groups", GroupsController.Delete)
	apiRoutes.PUT("/enable-peers", PeersController.Enable)
	apiRoutes.PUT("/enable-users", UserController.Enable)
	apiRoutes.POST("/execute", SystemController.Execute)
	apiRoutes.GET("/geofile", GeoController.GeoFile)
	apiRoutes.POST("/heartbeat", SystemController.HeartBeat)
	apiRoutes.GET("/id-change-support", SystemController.IDChangeSupportGet)
	apiRoutes.PUT("/id-change-support", SystemController.IDChangeSupportPut)
	apiRoutes.POST("/invite", UserController.Invite)
	apiRoutes.GET("/notices", SystemController.Notices)
	apiRoutes.POST("/reg", UserController.Reg)
	apiRoutes.POST("/reinvite", UserController.ReInvite)
	apiRoutes.GET("/resend/:username", UserController.Resend)
	apiRoutes.GET("/user-list", UserController.UserList)

	apiRoutesWithAuth := route.Group("/api").Use(Auth.TokenAuth())
	apiRoutesWithAuth.POST("/logout", UserController.Logout)
	apiRoutesWithAuth.POST("/currentUser", UserController.CurrentUser)
	apiRoutesWithAuth.POST("/initialState", SystemController.InitialState)

	AddressBookRoutes := apiRoutes.Group("/ab").Use(Auth.TokenAuth())
	AddressBookRoutes.GET("", AddressBookController.Get)
	AddressBookRoutes.POST("", AddressBookController.Update)
	AddressBookRoutes.DELETE("/bin/:id", AddressBookController.BinDelete)
	AddressBookRoutes.DELETE("/bin/empty/:id", AddressBookController.BinEmpty)
	AddressBookRoutes.PUT("/bin/restore/:id", AddressBookController.BinRestore)
	AddressBookRoutes.POST("/get", AddressBookController.Get)
	AddressBookRoutes.DELETE("/peer/:id", AddressBookController.PeerDelete)
	AddressBookRoutes.POST("/peer/add/:id", AddressBookController.PeerAdd)
	AddressBookRoutes.PUT("/peer/update/:id", AddressBookController.PeerUpdate)
	AddressBookRoutes.GET("/peers", AddressBookController.PeersGet)
	AddressBookRoutes.POST("/personal", AddressBookController.Personal)
	AddressBookRoutes.POST("/rule", AddressBookController.RulePost)
	AddressBookRoutes.PATCH("/rule", AddressBookController.RulePatch)
	AddressBookRoutes.POST("/rules", AddressBookController.RulesPost)
	AddressBookRoutes.DELETE("/rules", AddressBookController.RulesDelete)
	AddressBookRoutes.POST("/settings", AddressBookController.Settings)
	AddressBookRoutes.DELETE("/shared", AddressBookController.SharedDelete)
	AddressBookRoutes.POST("/shared/profiles", AddressBookController.SharedProfilesPost)
	AddressBookRoutes.POST("/shared/profile/:id", AddressBookController.SharedProfilePost)
	AddressBookRoutes.POST("/shared/add", AddressBookController.SharedProfileAdd)
	AddressBookRoutes.PUT("/shared/update/profile", AddressBookController.SharedUpdateProfile)
	AddressBookRoutes.DELETE("/tag/:id", AddressBookController.TagDelete)
	AddressBookRoutes.POST("/tag/add/:id", AddressBookController.TagAdd)
	AddressBookRoutes.PUT("/tag/update/:id", AddressBookController.TagUpdate)
	AddressBookRoutes.POST("/tags/:id", AddressBookController.TagsPost)

	AuditRoutes := apiRoutes.Group("/audit")
	AuditRoutes.GET("", AuditController.Audit).Use(Auth.TokenAuth())
	AuditRoutes.POST("", AuditController.Audit).Use(Auth.TokenAuth())
	AuditRoutes.PUT("", AuditController.AuditPut)
	AuditRoutes.POST("/disconnect", AuditController.Disconnect)

	AuditsRoutes := apiRoutes.Group("/audits")
	AuditsRoutes.GET("/alarm", AuditsController.Alarm)
	AuditsRoutes.GET("/conn", AuditsController.Conn)
	AuditsRoutes.GET("/console", AuditsController.Console)
	AuditsRoutes.GET("/file", AuditsController.File)

	CrossGroupRoutes := apiRoutes.Group("/cross-group")
	CrossGroupRoutes.POST("/create-rule", CrossGroupController.Create)
	CrossGroupRoutes.PUT("/delete-rule", CrossGroupController.Delete)
	CrossGroupRoutes.POST("/rules", CrossGroupController.RulesPost)

	CustomClientsRoutes := apiRoutes.Group("/custom-clients")
	CustomClientsRoutes.GET("", CustomClientsController.List)
	CustomClientsRoutes.POST("", CustomClientsController.Create)
	CustomClientsRoutes.GET("/:id", CustomClientsController.View)
	CustomClientsRoutes.PATCH("/:id", CustomClientsController.View)
	CustomClientsRoutes.DELETE("/:id", CustomClientsController.Delete)
	CustomClientsRoutes.POST("/:me/:le", CustomClientsController.Post)

	DevicesRoutes := apiRoutes.Group("/devices")
	DevicesRoutes.DELETE("/:id", DevicesController.Delete)

	GeoRoutes := apiRoutes.Group("/geo")
	GeoRoutes.GET("", GeoController.View)
	GeoRoutes.PUT("", GeoController.Update)
	GeoRoutes.DELETE("", GeoController.Delete)

	GroupRoutes := apiRoutes.Group("/group")
	GroupRoutes.POST("", GroupController.GroupPost)
	GroupRoutes.PUT("", GroupController.GroupPut)

	GroupsRoutes := apiRoutes.Group("/groups")
	GroupsRoutes.GET("", GroupsController.Get)

	KeypairRoutes := apiRoutes.Group("/keypair")
	KeypairRoutes.GET("", KeypairController.KeypairGet)
	KeypairRoutes.POST("", KeypairController.KeypairPost)
	KeypairRoutes.PUT("", KeypairController.KeypairPut)

	NamesRoutes := apiRoutes.Group("/names")
	NamesRoutes.GET("", NamesController.Get)

	PeerRoutes := apiRoutes.Group("/peer")
	PeerRoutes.PUT("", PeerController.Put)

	PeersRoutes := apiRoutes.Group("/peers")
	PeersRoutes.GET("", PeersController.Get)

	SessionRoutes := apiRoutes.Group("/session")
	SessionRoutes.GET("/exp_time", SessionController.ExpTimeGet)
	SessionRoutes.PUT("/exp_time", SessionController.ExpTimePut)

	SettingsRoutes := apiRoutes.Group("/settings")
	SettingsRoutes.PUT("", SettingsController.Put)

	SMTPRoutes := apiRoutes.Group("/smtp")
	SMTPRoutes.POST("/validate", SMTPController.Validate)

	StrategiesRoutes := apiRoutes.Group("/strategies")
	StrategiesRoutes.GET("", StrategiesController.Get)

	StrategyRoutes := apiRoutes.Group("/strategy")
	StrategyRoutes.GET("", StrategyController.Get)
	StrategyRoutes.PUT("", StrategyController.Put)
	StrategyRoutes.POST("/assign", StrategyController.Assign)
	StrategyRoutes.PUT("/op", StrategyController.Op)
	StrategyRoutes.GET("/options", StrategyController.Options)
	StrategyRoutes.PUT("/status", StrategyController.Status)

	TeamRoutes := apiRoutes.Group("/team")
	TeamRoutes.GET("/info", TeamController.InfoGet)
	TeamRoutes.PUT("/info", TeamController.InfoPut)

	TokensRoutes := apiRoutes.Group("/tokens")
	TokensRoutes.GET("", TokensController.Get)
	TokensRoutes.POST("", TokensController.Post)
	TokensRoutes.DELETE("/:token", TokensController.Delete)

	OIDCRoutes := apiRoutes.Group("/oidc")
	OIDCRoutes.POST("/auth", OIDCController.Auth)
	OIDCRoutes.GET("/callback", OIDCController.Callback)
	OIDCRoutes.GET("/settings", OIDCController.SettingsGet)
	OIDCRoutes.PUT("/settings", OIDCController.SettingsPut)
	OIDCRoutes.DELETE("/settings/:id", OIDCController.SettingsDelete)
	OIDCRoutes.POST("/settings/reorder", OIDCController.SettingsReorder)

	UserRoutes := apiRoutes.Group("/user")
	UserRoutes.POST("", UserController.Post)
	UserRoutes.PUT("", UserController.Put)
	UserRoutes.DELETE("/:id", UserController.Delete)
	UserRoutes.GET("/info", UserController.Info)
	UserRoutes.GET("/sessions", UserController.Sessions)
	UserRoutes.GET("/tfa/totp/backup-codes", UserController.TfaTotpBackupCodes)
	UserRoutes.PUT("/tfa/totp/disable", UserController.TfaTotpDisable)
	UserRoutes.POST("/tfa/totp/enable", UserController.TfaTotpEnable)
	UserRoutes.POST("/tfa/totp/verify", UserController.TfaTotpVerify)
	UserRoutes.PUT("/delete-session", UserController.DeleteSession)
	UserRoutes.PUT("/unsubscribe-email", UserController.UnsubscribeEmail)
	UserRoutes.POST("/update-verify", UserController.UpdateVerify)
	UserRoutes.POST("/verification-code", UserController.VerificationCode)

	VerifyRoutes := apiRoutes.Group("/verify")
	VerifyRoutes.POST("/user", VerifyController.User)
	VerifyRoutes.POST("/update", VerifyController.Update)
}
