package router

import (
	"IMCHAT/docs"
	"IMCHAT/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() *gin.Engine {
	r := gin.Default()
	//静态HTML资源

	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.StaticFile("/favicon.ico", "asset/images/favicon.ico")
	r.Static("/asset", "asset/")
	r.LoadHTMLGlob("views/**/*")
	//首页
	r.GET("/", service.GetIndex)
	r.GET("/index", service.GetIndex)
	r.GET("/register", service.Register)
	r.GET("/toChat", service.ToChat)
	r.GET("/chat", service.Chat)
	r.POST("/searchFriends", service.SearchFriends)

	//用户模块
	r.POST("/user/getUserList", service.GetUserList)
	r.POST("/user/createUser", service.CreateUser)
	r.POST("/user/deleteUser", service.DeleteUser)
	r.POST("/user/updateUser", service.UpdateUser)
	r.POST("/user/findUserByNameAndPwd", service.FindUserByNameAndPwd)
	r.POST("/user/find", service.FindByID)

	//发送消息
	r.GET("/user/sendMsg", service.SendMsg)
	r.GET("/user/SendUserMsg", service.SendUserMsg)

	//添加好友
	r.POST("/contact/addfriend", service.AddFriend)

	//文件中转
	r.POST("/attach/upload", service.Upload)
	//redis读取
	r.POST("/user/redisMsg", service.RedisMsg)

	//创建群
	r.POST("/contact/createCommunity", service.CreateCommunity)
	//群列表
	r.POST("/contact/loadcommunity", service.LoadCommunity)
	r.POST("/contact/joinGroup", service.JoinGroups)
	return r
}
