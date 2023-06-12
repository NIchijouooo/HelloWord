package models

import (
	"gateway/models"

	"gorm.io/gorm"
)

// 定义字典类型管理的存储库
type DictTypeRepository struct {
	db *gorm.DB
}

func NewDictTypeRepository() *DictTypeRepository {
	return &DictTypeRepository{db: models.DB}
}

// 新增字典类型
func (r *DictTypeRepository) Create(dictType *models.DictType) error {
	return r.db.Create(dictType).Error
}

// 修改字典类型
func (r *DictTypeRepository) Update(dictType *models.DictType) error {
	return r.db.Save(dictType).Error
}

// 删除字典类型
func (r *DictTypeRepository) Delete(dictId int) error {
	return r.db.Delete(&models.DictType{}, dictId).Error
}

// 获取所有字典类型
func (r *DictTypeRepository) GetAll(dictName, createTimeStart, createTimeEnd string, page, pageSize int) ([]models.DictType, int64, error) {
	var (
		dictTypeList []models.DictType
		total        int64
	)

	query := r.db.Model(&models.DictType{})
	if dictName != "" {
		query = query.Where("dict_name = ?", dictName)
	}
	if createTimeStart != "" {
		query = query.Where("created_at >= ?", createTimeStart)
	}
	if createTimeEnd != "" {
		query = query.Where("created_at <= ?", createTimeEnd)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&dictTypeList).Error; err != nil {
		return nil, 0, err
	}

	// paginator, err := paging.NewOffsetPaginator(query, uint64(offset), uint64(limit))
	// total = paginator.TotalEntries()
	// if err := paginator.Page(); err != nil {
	// 	return nil, 0, err
	// }
	return dictTypeList, total, nil
}

// 获取单个字典类型
func (r *DictTypeRepository) GetById(dictId int) (models.DictType, error) {
	var dictType models.DictType
	err := r.db.First(&dictType, dictId).Error
	return dictType, err
}
