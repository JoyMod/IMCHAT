package models

import (
	"IMCHAT/utils"
	"fmt"
	"gorm.io/gorm"
	"time"
)

//创建数据库对象表

type UserBasic struct {
	gorm.Model
	Name          string
	PassWord      string
	Phone         string `valid:"matches(^1[3-9]{1}\\d{9}$)"`
	Email         string `valid:"email"`
	Identity      string
	ClientIP      string
	ClientPort    string
	Salt          string    //随机字段
	LoginTime     time.Time //登录时间
	HeartbeatTime time.Time //心跳时间
	LoginOutTime  time.Time //下线时间
	IsLoginOut    bool      //登录状态
	DeviceInfo    string    //设备信息
}

//工厂模式

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10)
	utils.DB.Find(&data)
	for _, v := range data {
		fmt.Println(v)
	}
	return data
}

func FindUserByNameAndPwd(name string, pwd string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name=? and pass_word=?", name, pwd).First(&user)

	//token加密
	str := fmt.Sprintf("%d", time.Now().Unix())
	temp := utils.MD5Encode(str)
	utils.DB.Model(&user).Where("id = ?", user.ID).Update("identity", temp)
	return user
}

func FindUserByName(name string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name=?", name).First(&user)
	return user
}
func FindUserByPhone(phone string) *gorm.DB {
	user := UserBasic{}
	return utils.DB.Where("phone=?", phone).First(&user)
}

func FindUserByEmail(email string) *gorm.DB {
	user := UserBasic{}
	return utils.DB.Where("email=?", email).First(&user)
}

func CreateUser(user UserBasic) *gorm.DB {
	return utils.DB.Create(&user)
}

func DeleteUser(user UserBasic) *gorm.DB {
	return utils.DB.Delete(&user)
}

func UpdateUser(user UserBasic) *gorm.DB {
	return utils.DB.Model(&user).Updates(UserBasic{Name: user.Name, PassWord: user.PassWord, Phone: user.Phone, Email: user.Email})
}
func FindByID(id uint) UserBasic {
	user := UserBasic{}
	utils.DB.Where("id=?", id).First(&user)
	return user
}
