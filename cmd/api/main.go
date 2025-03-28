package main

import (
	"Vnollx/cmd/api/handler"
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
}
func RegisterGroup(r *gin.Engine) *gin.Engine {
	user := r.Group("/user")
	{
		user.GET("/captcha", handler.GetCaptcha)
		user.POST("/login", handler.Login)
		user.POST("/register", handler.Register)
		user.POST("/userinfo", middlerware.JWTAuthMiddleware(), handler.GetUserinfo)
		user.POST("/update", middlerware.JWTAuthMiddleware(), handler.UpdateUserInfo)
	}
	file := r.Group("/file")
	file.Use(middlerware.JWTAuthMiddleware())
	{
		file.POST("/upload", handler.UploadFile)
		file.POST("/delete", handler.DeleteFile)
		file.POST("/update", handler.UpdateFileInfo)
		file.POST("/enquery", handler.GetFileInfo)
		file.POST("/enquerys", handler.GetFilesInfo)
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
	return r
}
func main() {
	r := gin.Default()
	r.LoadHTMLGlob("D:/GolandProgram/Vnollx/web/*.html")
	LoadHTML(r)
	RegisterGroup(r)
	err := r.Run(apiServerAddr)
	if err != nil {
		log.Println(err)
		panic(err)
	}
}
