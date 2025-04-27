package models

import (
	"gorm.io/gorm"
	"time"
)

type UserBasic struct {
	ID        uint           `gorm:"primarykey;" json:"id"`
	CreatedAt MyTime         `gorm:"type:timestamp;" json:"created_at"`
	UpdatedAt MyTime         `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;" json:"deleted_at"`
	Identity  string         `gorm:"column:identity;type:varchar(36);" json:"identity"` // 用户的唯一标识
	Name      string         `gorm:"column:name;type:varchar(100);" json:"name"`        // 用户名
	Password  string         `gorm:"column:password;type:varchar(32);" json:"-"`        // 密码
	Phone     string         `gorm:"column:phone;type:varchar(20);" json:"phone"`       // 手机号
	Mail      string         `gorm:"column:mail;type:varchar(100);" json:"mail"`        // 邮箱
	PassNum   int64          `gorm:"column:pass_num;type:int(11);" json:"pass_num"`     // 通过的次数
	SubmitNum int64          `gorm:"column:submit_num;type:int(11);" json:"submit_num"` // 提交次数
	IsAdmin   int            `gorm:"column:is_admin;type:tinyint(1);" json:"is_admin"`  // 是否是管理员【0-否，1-是】
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}

func GetUserDetail(identity string) (data *UserBasic, err error) {
	data = new(UserBasic)
	err = DB.Omit("password").Where("identity = ? ", identity).Find(&data).Error
	return data, err
}

func GetUserOne(username, password string) (data *UserBasic, err error) {
	data = new(UserBasic)
	err = DB.Where("name = ? AND password = ? ", username, password).First(&data).Error
	return data, err
}

func UserCreate(identity, name, password, phone, mail string, is_admin int) (err error) {
	data := &UserBasic{
		Identity:  identity,
		Name:      name,
		Password:  password,
		Phone:     phone,
		Mail:      mail,
		CreatedAt: MyTime(time.Now()),
		UpdatedAt: MyTime(time.Now()),
		IsAdmin:   is_admin,
	}
	err = DB.Create(data).Error
	return err
}

func GetRankList(page, size int) (list []*UserBasic, count int64, err error) {
	list = make([]*UserBasic, 0)
	err = DB.Model(new(UserBasic)).Count(&count).Order("pass_num DESC, submit_num ASC").
		Offset(page).Limit(size).Find(&list).Error
	return list, count, err
}
