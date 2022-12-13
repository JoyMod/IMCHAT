package models

import (
	"IMCHAT/utils"
	"gorm.io/gorm"
)

//人员关系

type Contact struct {
	gorm.Model
	OwnerId  uint   //谁的关系信息
	TargetID uint   //对应谁
	Type     int    //对应类型
	Desc     string //保留字段
}

func (table *Contact) TableName() string {
	return "contact"
}

func SearchFriend(userId uint) []UserBasic {
	cons := make([]Contact, 0)
	obj := make([]uint64, 0)
	utils.DB.Where("owner_id = ? and type=1", userId).Find(&cons)

	for _, v := range cons {
		obj = append(obj, uint64(v.TargetID))
	}
	//查询好友列表
	users := make([]UserBasic, 0)
	utils.DB.Where("id in ?", obj).Find(&users)
	return users
}

func AddFriend(userId uint, targetName string) (int, string) {
	if targetName != "" {
		targetUser := FindUserByName(targetName)
		if targetUser.Salt != "" {
			if targetUser.ID == userId {
				return -1, "不能加自己"
			}
			contact0 := Contact{}
			utils.DB.Where("owner_id =?  and target_id =? and type=1", userId, targetUser.ID).Find(&contact0)
			if contact0.ID != 0 {
				return -1, "不能重复添加"
			}
			tx := utils.DB.Begin()
			//事务一旦开始，不论什么异常最终都会 Rollback
			defer func() {
				if r := recover(); r != nil {
					tx.Rollback()
				}
			}()
			contact := Contact{}
			contact.OwnerId = userId
			contact.TargetID = targetUser.ID
			contact.Type = 1
			if err := utils.DB.Create(&contact).Error; err != nil {
				tx.Rollback()
				return -1, "添加好友失败"
			}
			contact1 := Contact{}
			contact1.OwnerId = targetUser.ID
			contact1.TargetID = userId
			contact1.Type = 1
			if err := utils.DB.Create(&contact1).Error; err != nil {
				tx.Rollback()
				return -1, "添加好友失败"
			}
			tx.Commit()
			return 0, "添加好友成功"
		}
		return -1, "没有找到此用户"
	}
	return -1, "好友ID不能为空"
}

func SearchUserByGroupId(communityId uint) []uint {
	contacts := make([]Contact, 0)
	objIds := make([]uint, 0)
	utils.DB.Where("target_id = ? and type=2", communityId).Find(&contacts)
	for _, v := range contacts {
		objIds = append(objIds, uint(v.OwnerId))
	}
	return objIds
}
