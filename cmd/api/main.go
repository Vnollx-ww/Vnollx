package main

import (
	"Vnollx/cmd/api/handler"
	confg "Vnollx/config"
	"Vnollx/dal/db"
	"Vnollx/pkg/middlerware"
	"Vnollx/pkg/viper"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

var (
	config        = viper.Init("api")
	apiServerAddr = fmt.Sprintf("%s:%d", config.Viper.GetString("server.host"), config.Viper.GetInt("server.port"))
)

/*
	func LoadHTML(r *gin.Engine) {
		r.GET("/", func(c *gin.Context) {
			c.HTML(200, "index.html", nil) // 渲染 Web/index.html
		})
		r.GET("/login", func(c *gin.Context) {
			c.HTML(200, "login.html", nil) // 渲染 Web/index.html
		})
		r.GET("/register", func(c *gin.Context) {
			c.HTML(200, "register.html", nil) // 渲染 Web/index.html
		})
		r.GET("/homepage", func(c *gin.Context) {
			c.HTML(200, "homepage.html", nil) // 渲染 Web/index.html
		})
		r.GET("/info", func(c *gin.Context) {
			c.HTML(200, "info.html", nil)
		})
		r.GET("/allfile", func(c *gin.Context) {
			c.HTML(200, "allfile.html", nil)
		})
		r.GET("/help", func(c *gin.Context) {
			c.HTML(200, "help.html", nil)
		})
		r.GET("/setting", func(c *gin.Context) {
			c.HTML(200, "setting.html", nil)
		})
		r.GET("/category", func(c *gin.Context) {
			c.HTML(200, "category.html", nil)
		})
		r.GET("/share", func(c *gin.Context) {
			c.HTML(200, "share.html", nil)
		})
		r.GET("/notification", func(c *gin.Context) {
			c.HTML(200., "notification.html", nil)
		})
		r.GET("/friend", func(c *gin.Context) {
			c.HTML(200, "friend.html", nil)
		})
	}
*/
func RegisterGroup(r *gin.Engine) *gin.Engine {
	r.Use(middlerware.RateLimitMiddleware())
	user := r.Group("/user")
	{
		user.POST("/login", handler.Login)
		user.POST("/register", handler.Register)
		user.POST("/userinfo", middlerware.JWTAuthMiddleware(), handler.GetUserinfo)
		user.POST("/updateinfo", middlerware.JWTAuthMiddleware(), handler.UpdateUserInfo)
		user.POST("/updatepassword", middlerware.JWTAuthMiddleware(), handler.UpdatePassword)
		user.POST("/loginbycode", handler.UserLoginByCode)
		user.POST("/generatecaptcha", handler.GenerateCaptcha)
		user.POST("/addfriend", middlerware.JWTAuthMiddleware(), handler.AddFriend)
		user.POST("/deletefriend", middlerware.JWTAuthMiddleware(), handler.DeleteFriend)
		user.POST("/sendmessage", middlerware.JWTAuthMiddleware(), handler.SendMessage)
		user.POST("/getfriends", middlerware.JWTAuthMiddleware(), handler.GetFriendList)
		user.POST("/getmessages", middlerware.JWTAuthMiddleware(), handler.GetMessageList)
		user.POST("/deletemessage", middlerware.JWTAuthMiddleware(), handler.DeleteMessage)
		user.POST("/search", middlerware.JWTAuthMiddleware(), handler.GetUsersByName)
		user.GET("/ws", handler.WebSocketConnections)
		user.POST("/isfriend", middlerware.JWTAuthMiddleware(), handler.IsFriend) //后续请求记得判断是否需要加中间件
	}
	file := r.Group("/file")
	file.Use(middlerware.JWTAuthMiddleware())
	{
		file.POST("/upload", handler.UploadFile)
		file.POST("/delete", handler.DeleteFile)
		file.POST("/update", handler.UpdateFileInfo)
		file.POST("/enquery", handler.GetFileInfo)
		file.POST("/enquerys", handler.GetFilesInfo)
		file.POST("/getallfiles", handler.GetAllFiles)
		file.POST("/getvideofiles", handler.GetVideoFiles)
		file.POST("/getmusicfiles", handler.GetMusicFiles)
		file.POST("/getpicturefiles", handler.GetPictureFiles)
		file.POST("/getdocumentfiles", handler.GetDocumentFiles)
		file.POST("/getotherfiles", handler.GetOtherFiles)
		file.POST("/search", handler.GetFilesByName)
		file.POST("/createshare", handler.CreateShare)
		file.POST("/getfilebycode", handler.GetFileByCode)
	}
	folder := r.Group("/folder")
	folder.Use(middlerware.JWTAuthMiddleware())
	{
		folder.POST("/create", handler.CreateFolder)
		folder.POST("/delete", handler.DeleteFolder)
		folder.POST("/update", handler.UpdateFolderInfo)
		folder.POST("enquery", handler.GetFolderInfo)
		folder.POST("/enquerys", handler.GetFoldersInfo)
	}
	notification := r.Group("/notification")
	notification.Use(middlerware.JWTAuthMiddleware())
	{
		notification.POST("delete", handler.DeleteNotification)
		notification.POST("/getnotifications", handler.GetNotificationsByUserId)
		notification.POST("/sendnotification", handler.SendNotification)
	}
	return r
}
func main() {
	r := gin.Default()
	//r.LoadHTMLGlob("/app/web/*.html")
	//r.LoadHTMLGlob("D:/GolandProgram/Vnollx/web/*.html")
	//LoadHTML(r)
	r.Use(confg.Cors())
	RegisterGroup(r)
	go db.CleanExpiredShares()
	go db.SendNotificationIfSpaceNotEnough()
	err := r.Run(apiServerAddr)
	if err != nil {
		log.Println(err)
		panic(err)
	}
}
