package repositories

import (
	"fmt"
	"gateway/models"

	"gorm.io/gorm"
)

// 定义字典类型管理的存储库
type DictDataRepository struct {
	db *gorm.DB
}

func NewDictDataRepository() *DictDataRepository {
	return &DictDataRepository{db: models.DB}
}

// 新增字典类型
func (r *DictDataRepository) Create(dictData *models.DictData) error {
	return r.db.Create(dictData).Error
}

// 修改字典类型
func (r *DictDataRepository) Update(dictData *models.DictData) error {
	return r.db.Save(dictData).Error
}

// 删除字典类型
func (r *DictDataRepository) Delete(dictCode int) error {
	return r.db.Delete(&models.DictData{}, dictCode).Error
}

// 获取所有字典类型
func (r *DictDataRepository) GetAll(dictLabel, dictType string, page, pageSize int) ([]models.DictData, int64, error) {
	var (
		dictDataList []models.DictData
		total        int64
	)
	fmt.Println(dictLabel)
	fmt.Println(dictType)
	fmt.Println(page)
	fmt.Println(pageSize)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	query := r.db.Model(&models.DictData{})
	if dictLabel != "" {
		query = query.Where("dict_label = ?", dictLabel)
	}
	if dictType != "" {
		query = query.Where("dict_type = ?", dictType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	sql := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&dictDataList).Statement.SQL.String()
	fmt.Println(sql)

	return dictDataList, total, nil
}

// 获取单个字典类型
func (r *DictDataRepository) GetById(dictId int) (models.DictData, error) {
	fmt.Println(dictId)
	var dictData models.DictData
	err := r.db.First(&dictData, dictId).Error
	return dictData, err
}

// 获取字典类型下的所有字典数据
func (r *DictTypeRepository) GetDictDataListByDictTypeId(dictType string) ([]models.DictData, error) {
	var dictDataList []models.DictData
	err := r.db.Where("dict_type = ?", dictType).Find(&dictDataList).Error
	return dictDataList, err
}

// SelectDictValue /*根据字典类型和字典label获取字典信息*
func (r *DictDataRepository) SelectDictValue(dictType string, dictLabel string) (models.DictData, error) {
	var dictData models.DictData
	err := r.db.Where("dict_type = ?", dictType).Where("dict_label = ?", dictLabel).Find(&dictData).Error
	return dictData, err
}

func (r *DictDataRepository) GetDictDataByDictType(dictType string) ([]models.DictData, error) {
	var dictDataList []models.DictData
	err := r.db.Where("dict_type = ?", dictType).Find(&dictDataList).Error
	return dictDataList, err
}
