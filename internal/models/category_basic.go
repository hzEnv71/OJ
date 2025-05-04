package models

import (
	"gorm.io/gorm"
	"time"
)

type CategoryBasic struct {
	ID        uint           `gorm:"primarykey;" json:"id"`
	CreatedAt MyTime         `json:"created_at"`
	UpdatedAt MyTime         `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;" json:"deleted_at"`
	Identity  string         `gorm:"column:identity;type:varchar(36);" json:"identity"` // 分类的唯一标识
	Name      string         `gorm:"column:name;type:varchar(100);" json:"name"`        // 分类名称
	ParentId  int            `gorm:"column:parent_id;type:int(11);" json:"parent_id"`   // 父级ID
}

func (table *CategoryBasic) TableName() string {
	return "category_basic"
}

func GetCategoryList(keyword string, page int, size int) (data []*CategoryBasic, count int64, err error) {
	data = make([]*CategoryBasic, 0)
	err = DB.Model(new(CategoryBasic)).Where("name like ?", "%"+keyword+"%").
		Count(&count).Limit(size).Offset(page).Find(&data).Error
	return
}

func CategoryCreate(identity, name string, parentId int) (err error) {
	category := &CategoryBasic{
		Identity:  identity,
		Name:      name,
		ParentId:  parentId,
		CreatedAt: MyTime(time.Now()),
		UpdatedAt: MyTime(time.Now()),
	}
	err = DB.Create(category).Error
	return
}

func CategoryModify(identity, name string, parentId int) (err error) {
	category := &CategoryBasic{
		Identity:  identity,
		Name:      name,
		ParentId:  parentId,
		UpdatedAt: MyTime(time.Now()),
	}
	err = DB.Model(new(CategoryBasic)).Where("identity = ?", identity).Updates(category).Error
	return
}

func CategoryDelete(identity string) (err error) {
	err = DB.Where("identity = ?", identity).Delete(new(CategoryBasic)).Error
	return
}
