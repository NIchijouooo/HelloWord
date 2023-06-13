package models

// 定义字典类型表的模型
type DictType struct {
	DictId   int    `gorm:"primary_key;auto_increment;comment:'字典主键'" json:"dictId"` // 主键
	DictName string `gorm:"type:varchar(128);comment:'字典名称'" json:"dictName"`        // 字典名称
	DictType string `gorm:"type:varchar(64);comment:'字典类型'" json:"dictType"`         // 字典类型
	Status   string `gorm:"type:char(1);comment:'状态（0正常 1停用）'" json:"status"`        // 状态
}

func (u *DictType) TableName() string {
	return "sys_dict_type"
}
