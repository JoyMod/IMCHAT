package service

import (
	"IMCHAT/models"
	"IMCHAT/utils"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// GetUserList
// @Summary 所有用户
// @Tags 基础查询
// @Success 200 {string} json{"code,"message"}
// @Router /user/getUserList [get]

// GetUserList
// @Summary 所有用户
// @Tags 基础查询
// @Success 200 {string} json{"code,"message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()
	c.JSON(http.StatusOK, gin.H{"msg": data})
}

// CreateUser
// @Summary 新增用户
// @Tags 新增用户
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassword query string false "重复密码"
// @param phone query string false "电话号码"
// @param email query string false "电子邮件"
// @Success 200 {string} json{"code,"message"}
// @Router /user/createUser [get]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Name = c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	repassword := c.Request.FormValue("repassword")
	salt := fmt.Sprintf("%06d", rand.Int31())

	//用户名验证，如果存在，则不能注册
	dateName := models.FindUserByName(user.Name)
	if dateName.Name != "" {
		c.JSON(-1, gin.H{
			"code": -1, //0成功 -1失败
			"msg":  "用户名已存在！"})
		return
	}

	//非空验证
	if user.Name == "" || password == "" || repassword == "" {
		c.JSON(-1, gin.H{
			"code": -1, //0成功 -1失败
			"msg":  "用户名或密码不能为空！"})
		return
	}
	//密码验证
	if password != repassword {
		c.JSON(-1, gin.H{
			"code": -1, //0成功 -1失败
			"msg":  "两次密码不一致！"})
		return
	}
	//user.PassWord = password

	user.PassWord = utils.MakePasssword(password, salt)
	user.Salt = salt
	user.Phone = c.Request.FormValue("phone")
	user.Email = c.Request.FormValue("email")
	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(-1, gin.H{
			"code": -1, //0成功 -1失败
			"msg":  "电话或者电子邮件格式不对！"})
	} else {
		models.CreateUser(user)
		c.JSON(http.StatusOK, gin.H{
			"code": 0, //0成功 -1失败
			"msg":  "用户创建成功！"})
	}
}

// DeleteUser
// @Summary 删除用户
// @Tags 新增用户
// @param name query string false "id"
// @Router /user/deleteUser [get]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	models.DeleteUser(user)
	c.JSON(http.StatusOK, gin.H{
		"code": 0, //0成功 -1失败
		"msg":  "删除成功！",
	})
}

// UpdateUser
// @Summary 修改用户
// @Tags 新增用户
// @param id formData string false "id"
// @param name formData string false "name"
// @param password formData string false "password"
// @param phone formData string false "phone"
// @param email formData string false "email"
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.PassWord = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")

	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(-1, gin.H{
			"code": -1, //0成功 -1失败
			"msg":  "修改参数不匹配！",
		})
	} else {
		models.UpdateUser(user)
		c.JSON(http.StatusOK, gin.H{
			"code": 0, //0成功 -1失败
			"msg":  "修改成功！",
		})
	}
}

// FindUserByNameAndPwd
// @Summary 登录验证
// @Tags 新增用户
// @param name formData string false "用户名"
// @param password formData string false "密码"
// @Router /user/findUserByNameAndPwd [post]
func FindUserByNameAndPwd(c *gin.Context) {
	data := models.UserBasic{}

	name := c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	user := models.FindUserByName(name)
	if user.Name == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1, //0成功 -1失败
			"msg":  "信息不正确！",
			"data": data,
		})
		return
	}
	if !utils.ValidPasssword(password, user.Salt, user.PassWord) {
		c.JSON(http.StatusOK, gin.H{
			"code": -1, //0成功 -1失败
			"msg":  "信息不正确！",
			"data": data,
		})
		return
	}
	pwd := utils.MakePasssword(password, user.Salt)
	data = models.FindUserByNameAndPwd(name, pwd)

	c.JSON(http.StatusOK, gin.H{
		"code": 0, //0成功 -1失败
		"msg":  "登录成功！",
		"data": data,
	})
}

// 防止跨域伪造请求
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(c *gin.Context) {
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {

		}
	}(ws)
	MsgHandler(ws, c)
}

func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	for {
		msg, err := utils.Subscribe(c, utils.PublishKey)
		if err != nil {
			fmt.Println(err)
		}
		tm := time.Now().Format("2006-01-02 15:04:05")
		m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
		err = ws.WriteMessage(1, []byte(m))
		if err != nil {
			fmt.Println(err)
		}
	}

}
func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}

func SearchFriends(c *gin.Context) {
	w, _ := strconv.Atoi(c.Request.FormValue("userId"))
	b := models.SearchFriend(uint(w))
	utils.RespOKList(c.Writer, b, len(b))
}

func AddFriend(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userid"))
	targetName := c.Request.FormValue("targetName")
	code, msg := models.AddFriend(uint(userId), targetName)
	if code == 0 {
		utils.RespOK(c.Writer, code, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}
func RedisMsg(c *gin.Context) {
	userIdA, _ := strconv.Atoi(c.PostForm("userIdA"))
	userIdB, _ := strconv.Atoi(c.PostForm("userIdB"))
	start, _ := strconv.Atoi(c.PostForm("start"))
	end, _ := strconv.Atoi(c.PostForm("end"))
	isRev, _ := strconv.ParseBool(c.PostForm("isRev"))
	res := models.RedisMsg(int64(userIdA), int64(userIdB), int64(start), int64(end), isRev)
	utils.RespOKList(c.Writer, "ok", res)
}
func FindByID(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))

	//	name := c.Request.FormValue("name")
	data := models.FindByID(uint(userId))
	utils.RespOK(c.Writer, data, "ok")
}

func CreateCommunity(c *gin.Context) {
	ownerId, _ := strconv.Atoi(c.Request.FormValue("ownerId"))
	name := c.Request.FormValue("name")
	icon := c.Request.FormValue("icon")
	desc := c.Request.FormValue("desc")
	community := models.Community{}
	community.OwnerId = uint(ownerId)
	community.Name = name
	community.Img = icon
	community.Desc = desc
	code, msg := models.CreateCommunity(community)
	if code == 0 {
		utils.RespOK(c.Writer, code, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

//加载群列表

func LoadCommunity(c *gin.Context) {
	ownerId, _ := strconv.Atoi(c.Request.FormValue("ownerId"))
	//	name := c.Request.FormValue("name")
	data, msg := models.LoadCommunity(uint(ownerId))
	if len(data) != 0 {
		utils.RespList(c.Writer, 0, data, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}

//加群操作

func JoinGroups(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	comId := c.Request.FormValue("comId")

	//	name := c.Request.FormValue("name")
	data, msg := models.JoinGroup(uint(userId), comId)
	if data == 0 {
		utils.RespOK(c.Writer, data, msg)
	} else {
		utils.RespFail(c.Writer, msg)
	}
}
