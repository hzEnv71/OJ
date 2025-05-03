package models

import (
	"gorm.io/gorm"
)

type TestCase struct {
	ID              uint           `gorm:"primarykey;" json:"id"`
	CreatedAt       MyTime         `json:"created_at"`
	UpdatedAt       MyTime         `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index;" json:"deleted_at"`
	Identity        string         `gorm:"column:identity;type:varchar(255);" json:"identity"`
	ProblemIdentity string         `gorm:"column:problem_identity;type:varchar(255);" json:"problem_identity"`
	Input           string         `gorm:"column:input;type:text;" json:"input"`
	Output          string         `gorm:"column:output;type:text;" json:"output"`
}

func (table *TestCase) TableName() string {
	return "test_case"
}

func GetTestCase(identity string, page int, size int) (data []*TestCase, count int64, err error) {
	data = make([]*TestCase, 0)
	err = DB.Model(new(TestCase)).Where("problem_identity = ?", identity).Count(&count).Offset(page).Limit(size).Find(&data).Error
	return
}
