package models

import (
	"gorm.io/gorm"
	"oj/define"
	"oj/helper"
	"time"
)

type ProblemBasic struct {
	ID                uint               `gorm:"primarykey;" json:"id"`
	CreatedAt         MyTime             `json:"created_at"`
	UpdatedAt         MyTime             `json:"updated_at"`
	DeletedAt         gorm.DeletedAt     `gorm:"index;" json:"deleted_at"`
	Identity          string             `gorm:"column:identity;type:varchar(36);" json:"identity"`                  // 问题表的唯一标识
	ProblemCategories []*ProblemCategory `gorm:"foreignKey:problem_id;references:id" json:"problem_categories"`      // 关联问题分类表
	Title             string             `gorm:"column:title;type:varchar(255);" json:"title"`                       // 文章标题
	Content           string             `gorm:"column:content;type:text;" json:"content"`                           // 文章正文
	MaxRuntime        int                `gorm:"column:max_runtime;type:int(11);" json:"max_runtime"`                // 最大运行时长
	MaxMem            int                `gorm:"column:max_mem;type:int(11);" json:"max_mem"`                        // 最大运行内存
	TestCases         []*TestCase        `gorm:"foreignKey:problem_identity;references:identity;" json:"test_cases"` // 关联测试用例表
	PassNum           int64              `gorm:"column:pass_num;type:int(11);" json:"pass_num"`                      // 通过次数
	SubmitNum         int64              `gorm:"column:submit_num;type:int(11);" json:"submit_num"`                  // 提交次数
}

func (table *ProblemBasic) TableName() string {
	return "problem_basic"
}

func GetProblemList(keyword, categoryIdentity string, page, size int) (data []*ProblemBasic, count int64, err error) {
	data = make([]*ProblemBasic, 0)
	tx := DB.Model(new(ProblemBasic)).
		Distinct(`problem_basic.id`).
		Select(`problem_basic.id`, `problem_basic.identity`, `problem_basic.title`, `problem_basic.max_runtime`, `problem_basic.max_mem`, `problem_basic.pass_num`, `submit_num`, `problem_basic.created_at`, `problem_basic.updated_at`, `problem_basic.deleted_at`).
		Preload("ProblemCategories").
		Preload("ProblemCategories.CategoryBasic").
		Where("title like ? OR content like ? ", "%"+keyword+"%", "%"+keyword+"%")
	if categoryIdentity != "" {
		tx.Joins("RIGHT JOIN problem_category pc on pc.problem_id = problem_basic.id").
			Where("pc.category_id = (SELECT cb.id FROM category_basic cb WHERE cb.identity = ? )", categoryIdentity)
	}
	err = tx.Order("problem_basic.id DESC").Offset(page).Limit(size).Find(&data).Distinct(`problem_basic.id`).Count(&count).Error
	return data, count, err
}

func GetProblemDetail(identity string) (data *ProblemBasic, err error) {
	err = DB.Where("identity = ?", identity).
		Preload("ProblemCategories").Preload("ProblemCategories.CategoryBasic").
		First(&data).Error
	return data, err
}

func ProblemCreate(identity string, in *define.ProblemBasic) (err error) {
	data := &ProblemBasic{
		Identity:   identity,
		Title:      in.Title,
		Content:    in.Content,
		MaxRuntime: in.MaxRuntime,
		MaxMem:     in.MaxMem,
		CreatedAt:  MyTime(time.Now()),
		UpdatedAt:  MyTime(time.Now()),
	}
	// 处理分类
	problemCategories := make([]*ProblemCategory, 0)
	for _, id := range in.ProblemCategories {
		pc := &ProblemCategory{
			ProblemId:  data.ID,
			CategoryId: uint(id),
			CreatedAt:  MyTime(time.Now()),
			UpdatedAt:  MyTime(time.Now()),
		}
		problemCategories = append(problemCategories, pc)
	}
	data.ProblemCategories = problemCategories
	// 处理测试用例
	testCaseBasics := make([]*TestCase, 0)
	for _, v := range in.TestCases {
		// 举个例子 {"input":"1 2\n","output":"3\n"}
		testCaseBasic := &TestCase{
			Identity:        helper.GetUUID(),
			ProblemIdentity: identity,
			Input:           v.Input,
			Output:          v.Output,
			CreatedAt:       MyTime(time.Now()),
			UpdatedAt:       MyTime(time.Now()),
		}
		testCaseBasics = append(testCaseBasics, testCaseBasic)
	}
	data.TestCases = testCaseBasics
	// 创建问题
	err = DB.Create(data).Error
	return err
}

func ProblemModify(in *define.ProblemBasic) (err error) {
	var modify = func(tx *gorm.DB) error {
		// 问题基础信息保存 problem_basic
		problemBasic := &ProblemBasic{
			Identity:   in.Identity,
			Title:      in.Title,
			Content:    in.Content,
			MaxRuntime: in.MaxRuntime,
			MaxMem:     in.MaxMem,
			UpdatedAt:  MyTime(time.Now()),
		}
		// 更新问题
		err := tx.Where("identity = ?", in.Identity).Updates(problemBasic).Error
		if err != nil {
			return err
		}
		// 查询问题详情
		err = tx.Where("identity = ?", in.Identity).Find(problemBasic).Error
		if err != nil {
			return err
		}

		// 关联问题分类的更新
		// 1、删除已存在的关联关系
		err = tx.Where("problem_id = ?", problemBasic.ID).Delete(new(ProblemCategory)).Error
		if err != nil {
			return err
		}
		// 2、新增新的关联关系
		pcs := make([]*ProblemCategory, 0)
		for _, id := range in.ProblemCategories {
			pcs = append(pcs, &ProblemCategory{
				ProblemId:  problemBasic.ID,
				CategoryId: uint(id),
				CreatedAt:  MyTime(time.Now()),
				UpdatedAt:  MyTime(time.Now()),
			})
		}
		err = tx.Create(&pcs).Error
		if err != nil {
			return err
		}
		// 关联测试案例的更新
		// 1、删除已存在的关联关系
		err = tx.Where("problem_identity = ?", in.Identity).Delete(new(TestCase)).Error
		if err != nil {
			return err
		}
		// 2、增加新的关联关系
		tcs := make([]*TestCase, 0)
		for _, v := range in.TestCases {
			// 举个例子 {"input":"1 2\n","output":"3\n"}
			tc := &TestCase{
				Identity:        helper.GetUUID(),
				ProblemIdentity: in.Identity,
				Input:           v.Input,
				Output:          v.Output,
				CreatedAt:       MyTime(time.Now()),
				UpdatedAt:       MyTime(time.Now()),
			}
			tcs = append(tcs, tc)
		}
		err = tx.Create(tcs).Error
		if err != nil {
			return err
		}
		return nil
	}
	DB.Transaction(modify)
	return err
}
