package repositories

import (
	"gateway/models"
	"gorm.io/gorm"
)

// 定义字典类型管理的存储库
type ProjectInfoRepository struct {
	db *gorm.DB
}

func NewProjectInfoRepository() *ProjectInfoRepository {
	return &ProjectInfoRepository{db: models.DB}
}

// 新增字典类型
func (r *ProjectInfoRepository) Create(dictType *models.ProjectInfo) error {
	return r.db.Create(dictType).Error
}

// 修改字典类型
func (r *ProjectInfoRepository) Update(dictType *models.ProjectInfo) error {
	return r.db.Save(dictType).Error
}

// 删除字典类型
func (r *ProjectInfoRepository) Delete(dictId int) error {
	return r.db.Delete(&models.ProjectInfo{}, dictId).Error
}

// 获取所有字典类型
func (r *ProjectInfoRepository) GetAll(projectName, createTimeStart, createTimeEnd string) ([]models.ProjectInfo, error) {
	var dictTypeList []models.ProjectInfo
	query := r.db.Model(&models.ProjectInfo{})
	if projectName != "" {
		query = query.Where("name = ?", projectName)
	}
	if createTimeStart != "" {
		query = query.Where("created_at >= ?", createTimeStart)
	}
	if createTimeEnd != "" {
		query = query.Where("created_at <= ?", createTimeEnd)
	}
	if err := query.Find(&dictTypeList).Error; err != nil {
		return nil, err
	}
	err := r.db.Find(&dictTypeList).Error
	return dictTypeList, err
}

// 获取单个字典类型
func (r *ProjectInfoRepository) GetById(dictId int) (models.ProjectInfo, error) {
	var dictType models.ProjectInfo
	err := r.db.First(&dictType, dictId).Error
	return dictType, err
}
