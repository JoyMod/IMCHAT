package models

import "gorm.io/gorm"

//群信息

type GroupBasic struct {
	gorm.Model
	Name    string //群名称
	OwnerId uint   //群拥有者
	Icon    string
	Type    int //保留字段1
	Desc    string
}

func (table *GroupBasic) TableName() string {
	return "group_basic"
}
