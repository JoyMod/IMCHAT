package service

import (
	"IMCHAT/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"text/template"
)

// GetIndex
// @Tags 首页
// @Success 200 {string} welcome
// @Router /index [get]
func GetIndex(c *gin.Context) {
	ind, err := template.ParseFiles("index.html", "views/chat/head.html")
	if err != nil {
		panic(err)
	}
	err = ind.Execute(c.Writer, "index")
	if err != nil {
		return
	}
}

// Register
// @Tags 注册页面
// @Success 200 {string} welcome
// @Router /index [get]
func Register(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", nil)
}

func ToChat(c *gin.Context) {
	ind, err := template.ParseFiles("views/chat/index.html",
		"views/chat/head.html",
		"views/chat/foot.html",
		"views/chat/tabmenu.html",
		"views/chat/concat.html",
		"views/chat/group.html",
		"views/chat/profile.html",
		"views/chat/createcom.html",
		"views/chat/userinfo.html",
		"views/chat/main.html")
	if err != nil {
		panic(err)
	}
	userId, _ := strconv.Atoi(c.Query("userId"))
	token := c.Query("token")
	user := models.UserBasic{}
	user.ID = uint(userId)
	user.Identity = token
	err = ind.Execute(c.Writer, user)
	if err != nil {
		return
	}

}
func Chat(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}
