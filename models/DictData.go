package models

// 定义字典数据表的模型
type DictData struct {
	DictCode  int    `gorm:"primary_key;auto_increment;comment:'字典编码'" json:"dictCode"` // 主键
	DictSort  int    `gorm:"comment:'字典排序'" json:"dictSort"`                            // 排序
	DictLabel string `gorm:"type:varchar(128);comment:'字典标签'" json:"dictLabel"`         // 标签
	DictValue string `gorm:"type:varchar(255);comment:'字典键值'" json:"dictValue"`         // 键值
	DictType  string `gorm:"type:varchar(64);comment:'字典类型'" json:"dictType"`           // 类型
	CssClass  string `gorm:"type:varchar(128);comment:'样式属性（其他样式扩展）'" json:"cssClass"`
	ListClass string `gorm:"type:varchar(128);comment:'表格回显样式'" json:"listClass"` // 表格回显样式
	IsDefault string `gorm:"type:char(1);comment:'是否默认（Y是 N否）'" json:"isDefault"` // 是否默认
	Status    string `gorm:"type:char(1);comment:'状态（0正常 1停用）'" json:"status"`    // 状态
	Remark    string `gorm:"type:varchar(255);comment:'备注'" json:"remark"`        // 备注
}

func (u *DictData) TableName() string {
	return "sys_dict_data"
}
